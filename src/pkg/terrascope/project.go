package terrascope

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/sirupsen/logrus"
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
	sfm            *scopeFilterMatcher

	RootsDir string `hcl:"rootsDir"`
	Roots    map[string]*root

	*logrus.Entry
}

// ScopeType represents a single scope available to a project
type ScopeType struct {
	Name         string `hcl:"name"`
	Description  string `hcl:"description,optional"`
	DefaultValue string `hcl:"default,optional"`
	// Validations  []*ProjectScopeValidation `hcl:"validation,block"`
}

func (sc *ScopeType) String() string {
	return sc.Name
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

func (p *Project) scopeFilterMatcher() *scopeFilterMatcher {
	if p.sfm != nil {
		return p.sfm
	}
	p.sfm = newScopeFilterMatcher(p.compiledScopes, p.ScopeTypes, p.Logger)
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
	if err := p.assertRootDependenciesAcyclic(); err != nil {
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

	rootExec, err := newRootExecutor(root, scopes, p.scopeFilterMatcher(), chain, p.Logger)
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
