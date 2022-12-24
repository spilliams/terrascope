package terraboots

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spilliams/terraboots/internal/scopedata"
)

// GenerateScopeData builds a generator for new scope data, then executes it,
// saving the results in a file
func (p *Project) GenerateScopeData(input io.Reader, output io.Writer) error {
	if len(p.Scopes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.Scopes))
	for i, el := range p.Scopes {
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

// readScopeData reads all of the scope data known to the receiver
func (p *Project) readScopeData() error {
	if len(p.Scopes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.Scopes))
	for i, el := range p.Scopes {
		scopeTypes[i] = el.Name
	}

	filenames := make([]string, len(p.ScopeDataFiles))
	for i, filename := range p.ScopeDataFiles {
		filenames[i] = path.Join(p.projectDir(), filename)
	}
	reader := scopedata.NewReader(scopeTypes, filenames, p.Logger)
	var err error
	p.rootScopeValues, err = reader.Read()
	if err != nil {
		return err
	}

	return nil
}
