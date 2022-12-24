package terraboots

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// ProjectConfig represents the configuration file of a Terraboots project
type ProjectConfig struct {
	Project *Project `hcl:"terraboots,block"`
}

// Project represents a single Terraboots project, complete with scope types,
// scope data, and root configurations.
type Project struct {
	configFile string
	ID         string `hcl:"id,label"`

	Scopes    []*ProjectScope `hcl:"scope,block"`
	ScopeData []string        `hcl:"scopeData"`

	RootsDir string `hcl:"rootsDir"`
	Roots    map[string]*Root
}

// ProjectScope represents a single scope available to a project
type ProjectScope struct {
	Name         string `hcl:"name"`
	Description  string `hcl:"description,optional"`
	DefaultValue string `hcl:"default,optional"`
	// Validations  []*ProjectScopeValidation `hcl:"validation,block"`
}

// ProjectScopeValidation
// type ProjectScopeValidation struct {
// 	Condition    bool   `hcl:"condition"`
// 	ErrorMessage string `hcl:"error_message"`
// }

// ParseProject reads the given configuration file and parses it as a new
// Terraboots project
func ParseProject(cfgFile string) (*Project, error) {
	cfg := &ProjectConfig{}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return nil, err
	}
	err = hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Project.readScopeData()
	if err != nil {
		return nil, err
	}

	cfg.Project.configFile = cfgFile

	return cfg.Project, nil
}

// BuildRoot tells the receiver to build a root module
func (p *Project) BuildRoot(rootName string) (*Root, error) {
	root, ok := p.Roots[rootName]
	if !ok {
		var err error
		root, err = p.AddRoot(rootName)
		if err != nil {
			return nil, err
		}
	}

	// TODO
	return root, nil
}

// AddRoot tells the receiver to add a root module to its internal records.
// The `rootName` must be a directory name located in the receiver's `RootsDir`.
func (p *Project) AddRoot(rootName string) (*Root, error) {
	// look for named root
	rootDir := path.Join(p.RootsDir, rootName)
	_, err := os.Stat(rootDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("could not locate a root named '%s' in the roots directory '%s'", rootName, p.RootsDir)
	} else if err != nil {
		return nil, err
	}

	// look for terraboots file
	rootCfg := path.Join(rootDir, "terraboots.hcl")
	_, err = os.Stat(rootCfg)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("found a root named '%s' in the roots directory '%s', but it does not contain a terraboots.hcl configuration", rootName, p.RootsDir)
	} else if err != nil {
		return nil, err
	}

	root, err := ParseRoot(rootCfg)
	if err != nil {
		return nil, err
	}

	if p.Roots == nil {
		p.Roots = make(map[string]*Root)
	}
	p.Roots[root.ID] = root
	return root, nil
}