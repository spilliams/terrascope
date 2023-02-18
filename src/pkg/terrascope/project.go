package terrascope

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/sirupsen/logrus"
	"github.com/spilliams/terrascope/internal/generate"
	hclhelp "github.com/spilliams/terrascope/internal/hcl"
)

// ProjectConfig represents the configuration file of a Terrascope project
type ProjectConfig struct {
	Project *Project `hcl:"terrascope,block"`
}

// Project represents a single Terrascope project, complete with scope types,
// scope data, and root configurations.
type Project struct {
	configFile string
	ID         string `hcl:"id,label"`

	ScopeTypes     []*ScopeType `hcl:"scope,block"`
	ScopeDataFiles []string     `hcl:"scopeData"`
	compiledScopes CompiledScopes
	sfm            *scopeMatcher

	RootsDir string `hcl:"rootsDir"`
	Roots    map[string]*root

	*logrus.Entry
}

type scopeDataConfig struct {
	RootScopes []*NestedScope `hcl:"scope,block"`
}

// ProjectScopeValidation
// type ProjectScopeValidation struct {
// 	Condition    bool   `hcl:"condition"`
// 	ErrorMessage string `hcl:"error_message"`
// }

// ParseProject reads the given configuration file and parses it as a new
// Terrascope project
func ParseProject(cfgFile string, logger *logrus.Logger) (*Project, error) {
	cfg := &ProjectConfig{}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return nil, err
	}
	err = hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}

	project := cfg.Project
	project.configFile = cfgFile
	project.Entry = logger.WithField("prefix", "project")

	err = project.readScopeData()
	if err != nil {
		return nil, err
	}

	project.Debugf("Project has %d compiled scopes", project.compiledScopes.Len())
	for _, scope := range project.compiledScopes {
		project.Trace(scope.Address())
	}
	return project, nil
}

func (p *Project) projectDir() string {
	return path.Dir(p.configFile)
}

func (p *Project) scopeFilterMatcher() *scopeMatcher {
	if p.sfm != nil {
		return p.sfm
	}
	p.sfm = newScopeMatcher(p.compiledScopes, p.ScopeTypes, p.Logger)
	return p.sfm
}

// AddAllRoots searches the receiver's `RootsDir` for directories, and adds them
// all to the project as root configurations.
func (p *Project) AddAllRoots() error {
	files, err := ioutil.ReadDir(p.RootsDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			err := p.addRoot(file.Name())
			if err != nil {
				return err
			}
		}
	}

	// check for dependency-cycles
	rdc := &rootDependencyCalculator{roots: p.Roots}
	if err := rdc.assertRootDependenciesAcyclic(); err != nil {
		return err
	}

	return nil
}

// BuildRoot tells the receiver to build a root module, and returns a list of
// directories where the root was built to.
// If `chain` is `RootExecutorDependencyChainingUnknown`, this function will
// survey the user for a "none/one/all" choice pertaining to the root's
// dependencies.
func (p *Project) BuildRoot(rootName string, scopes []string, dryRun bool, chain RootExecutorDependencyChaining) ([]string, error) {
	// make sure the root exists
	root, ok := p.Roots[rootName]
	if !ok {
		return nil, fmt.Errorf("Root '%s' isn't loaded. Did you run `AddAllRoots`?", rootName)
	}
	root = p.Roots[rootName]
	rdc := &rootDependencyCalculator{
		roots: p.Roots,
		chain: chain,
	}
	rootExec, err := newRootExecutor(root, scopes, p.scopeFilterMatcher(), rdc, chain, p.Logger)
	if err != nil {
		return nil, err
	}

	return rootExec.Execute(BuildContext, dryRun)
}

// addRoot tells the receiver to add a root module to its internal records.
// The `rootName` must be a directory name located in the receiver's `RootsDir`.
func (p *Project) addRoot(rootName string) error {
	// look for named root
	rootDir := path.Join(p.RootsDir, rootName)
	_, err := os.Stat(rootDir)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not locate a root named '%s' in the roots directory '%s'", rootName, p.RootsDir)
	} else if err != nil {
		return err
	}
	p.Debugf("Adding root %s", rootDir)

	// look for terrascope file
	rootCfg := path.Join(rootDir, "terrascope.hcl")
	_, err = os.Stat(rootCfg)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("found a root named '%s' in the roots directory '%s', but it does not contain a terrascope.hcl configuration", rootName, p.RootsDir)
	} else if err != nil {
		return err
	}

	r, err := p.ParseRoot(rootCfg)
	if r == nil && err != nil {
		return err
	}

	if p.Roots == nil {
		p.Roots = make(map[string]*root)
	}
	if r != nil {
		p.Roots[r.name] = r
	}
	return nil
}

func pluralize(single, plural string, count int) string {
	if count == 1 {
		return single
	}
	return plural
}

// ParseRoot tells the receiver to parse a root module configuration file at the
// given path.
func (p *Project) ParseRoot(cfgFile string) (*root, error) {
	// partial decode first, because we don't know what scope or attributes
	// this config will use. We're just looking for the `root` block here.
	cfg := &struct {
		Root *root `hcl:"root,block"`
	}{}

	err := hclsimple.DecodeFile(cfgFile, hclhelp.DefaultContext(), cfg)
	r := cfg.Root
	// we purposefully ignore err until the end
	if r == nil {
		p.Warnf("Root detected at %s failed to decode. Does it have a complete terrascope.hcl file?", cfgFile)
		return nil, nil
	}

	r.filename = cfgFile
	r.name = path.Base(path.Dir(cfgFile))
	return r, err
}

// GenerateRoot creates and runs a new root generator, using the receiver's
// scope types.
func (p *Project) GenerateRoot(name string) error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terrascope `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.ScopeTypes))
	for i, el := range p.ScopeTypes {
		scopeTypes[i] = el.Name
	}

	return generate.Root(name, p.RootsDir, scopeTypes, p.Logger)
}

// GenerateScopeData builds a generator for new scope data, then executes it,
// saving the results in a file
func (p *Project) GenerateScopeData() error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terrascope `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.ScopeTypes))
	for i, el := range p.ScopeTypes {
		scopeTypes[i] = el.Name
	}

	// this file doesn't have to exist yet
	dataFilename := "data.hcl"
	if p.ScopeDataFiles != nil && len(p.ScopeDataFiles) > 0 {
		// TODO: which filename? a new one? and then update the project config with the new filename?
		dataFilename = p.ScopeDataFiles[0]
	}
	dataFilename = path.Join(p.projectDir(), dataFilename)

	return generate.Scope(scopeTypes, dataFilename, p.Logger)
}

// readScopeData reads all of the scope data known to the receiver
func (p *Project) readScopeData() error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terrascope `scope` block, then try this again", p.configFile)
	}
	if len(p.compiledScopes) > 0 {
		return nil
	}

	list := make([]*CompiledScope, 0)

	for _, filename := range p.ScopeDataFiles {
		filename := path.Join(p.projectDir(), filename)
		p.Debugf("Reading scope data file %s", filename)
		if !fileExists(filename) {
			continue
		}

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

	compiledScopes := CompiledScopes(list)
	compiledScopes = compiledScopes.Deduplicate()
	sort.Sort(compiledScopes)

	p.compiledScopes = compiledScopes
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
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

// GetCompiledScopes returns a list of the receiver's compiled scopes that match
// a given address. For more information on how a scope matches an address
// string, see `CompiledScope.Matches`.
func (p *Project) GetCompiledScopes(address string) (CompiledScopes, error) {
	if err := p.readScopeData(); err != nil {
		return nil, err
	}

	scopes := CompiledScopes{}
	filter, err := p.scopeFilterMatcher().makeFilter(address)
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

	filter, err := p.scopeFilterMatcher().makeFilter(address)
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
