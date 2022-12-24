package terraboots

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type ProjectConfig struct {
	Project *Project `hcl:"terraboots,block"`
}

type Project struct {
	configFile string
	ID         string   `hcl:"id,label"`
	RootsDir   string   `hcl:"rootsDir"`
	ScopeData  []string `hcl:"scopeData"`

	Scopes []*ProjectScope `hcl:"scope,block"`
	Roots  map[string]*Root
}

type ProjectScope struct {
	Name         string                    `hcl:"name"`
	Description  string                    `hcl:"description,optional"`
	DefaultValue string                    `hcl:"default,optional"`
	Validations  []*ProjectScopeValidation `hcl:"validation,block"`
}

type ProjectScopeValidation struct {
	Condition    bool   `hcl:"condition"`
	ErrorMessage string `hcl:"error_message"`
}

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

	// TODO: read scope data

	cfg.Project.configFile = cfgFile

	return cfg.Project, nil
}

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
