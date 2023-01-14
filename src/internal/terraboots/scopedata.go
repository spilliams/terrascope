package terraboots

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spilliams/terraboots/internal/hclhelp"
	"github.com/spilliams/terraboots/internal/scopedata"
)

// GenerateScopeData builds a generator for new scope data, then executes it,
// saving the results in a file
func (p *Project) GenerateScopeData(input io.Reader, output io.Writer) error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.ScopeTypes))
	for i, el := range p.ScopeTypes {
		scopeTypes[i] = el.Name
	}

	gen := scopedata.NewGenerator(scopeTypes, p.Logger)
	bytes, err := gen.Create(input, output)
	if err != nil {
		return err
	}

	// this file doesn't have to exist yet
	dataFilename := "data.hcl"
	if p.ScopeDataFiles != nil && len(p.ScopeDataFiles) > 0 {
		// TODO: which filename? a new one? and then update the project config with the new filename?
		dataFilename = p.ScopeDataFiles[0]
	}
	dataFilename = path.Join(p.projectDir(), dataFilename)

	file, err := os.OpenFile(dataFilename, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			p.Debug("file open errored, is ErrNotExist, creating file")
			file, err = os.Create(dataFilename)
			if err != nil {
				p.Debug("file create failed")
				return err
			}
		} else {
			p.Debug("file open errored, is not ErrNotExist, throwing")
			return err
		}
	} else {
		// err == nil means file was found
		p.Warnf("A file '%s' already exists! Overwrite? [Y/n]", dataFilename)
		scanner := bufio.NewScanner(input)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			p.Debug("scanner errored")
			return err
		}
		if len(scanner.Text()) != 0 {
			p.Debug("scanner returned text")
			if scanner.Text() != "y" && scanner.Text() != "Y" {
				p.Debugf("User does not want to overwrite, printing and exiting.")
				output.Write(bytes)
				return nil
			}
		}
	}
	_, err = file.Write(bytes)
	return err
}

type scopeDataConfig struct {
	RootScopes []*scopedata.NestedScope `hcl:"scope,block"`
}

// readScopeData reads all of the scope data known to the receiver
func (p *Project) readScopeData() error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}
	if len(p.compiledScopes) > 0 {
		return nil
	}

	list := make([]*scopedata.CompiledScope, 0)

	for _, filename := range p.ScopeDataFiles {
		filename := path.Join(p.projectDir(), filename)
		p.Debugf("Reading scope data file %s", filename)

		cfg := &scopeDataConfig{}
		err := hclsimple.DecodeFile(filename, nil, cfg)
		if err = handleDecodeNestedScopeError(err); err != nil {
			p.Warnf("error decoding scope data file %s", filename)
			return err
		}

		for _, rootScope := range cfg.RootScopes {
			list = append(list, rootScope.CompiledScopes(nil)...)
		}
	}

	compiledScopes := scopedata.CompiledScopes(list)
	compiledScopes = compiledScopes.Deduplicate()
	sort.Sort(compiledScopes)

	p.compiledScopes = compiledScopes
	return nil
}

// handleDecodeNestedScopeError takes diagnostics returned from a call to decode
// something into a NestedScope, and it removes the diagnostics that are false
// alarms.
// When dealing with the `remain` tag in a struct, gohcl will add a diagnostic
// that "Blocks are not allowed here". It's ok to ignore this type of diagnostic
// because the blocks are handled elsewhere in the gohcl Decode process.
// Errors that are not hcl Diagnostics, or that are other types of Diagnostic
// will be returned.
func handleDecodeNestedScopeError(err error) error {
	return hclhelp.DiagnosticsWithoutSummary(err, "Unexpected \"scope\" block")
}

func (p *Project) GetCompiledScopes(address string) (scopedata.CompiledScopes, error) {
	scopes := scopedata.CompiledScopes{}
	filter, err := p.makeScopeFilter(address)
	if err != nil {
		return nil, err
	}
	for _, scope := range p.compiledScopes {
		ok, err := scope.Matches(filter)
		if err != nil {
			return nil, err
		}
		if ok {
			scopes = append(scopes, scope)
		}
	}
	return scopes, nil
}

// IsScopeValue checks the given address against the receiver's list of known
// scope values. May return an error if the receiver can't read its scope
// values.
func (p *Project) IsScopeValue(address string) (bool, error) {
	if err := p.readScopeData(); err != nil {
		return false, err
	}

	filter, err := p.makeScopeFilter(address)
	if err != nil {
		return false, err
	}

	for _, scope := range p.compiledScopes {
		matches, err := scope.Matches(filter)
		if err != nil {
			return false, err
		}
		p.Tracef("%s matches? %v", scope, matches)
		if matches {
			return true, nil
		}
	}
	return false, nil
}

// address could be types & values interleaved, or just values
func (p *Project) makeScopeFilter(address string) (map[string]string, error) {
	p.Debugf("makeScopeFilter %s", address)

	m := make(map[string]string)
	parts := strings.Split(address, ".")

	if len(parts)%2 == 0 {
		isCollated := true
		i := 0
		for i*2 < len(parts) {
			if parts[i*2] != p.ScopeTypes[i].Name {
				isCollated = false
				break
			}
			i++
		}

		if isCollated {
			newParts := make([]string, 0)
			for i := 1; i < len(parts); i += 2 {
				newParts = append(newParts, parts[i])
			}
			parts = newParts
		}
	}
	p.Debugf("  parts after decollation: %v", parts)

	if len(parts) > len(p.ScopeTypes) {
		return nil, fmt.Errorf("scope address %s is too long to be mapped to scope types %v", address, p.ScopeTypes)
	}
	for i, v := range parts {
		// some special regex translating
		if v == "*" {
			v = ".*"
		}
		m[p.ScopeTypes[i].Name] = v
	}
	p.Debugf("  mapping %+v", m)
	return m, nil
}
