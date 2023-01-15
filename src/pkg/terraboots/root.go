package terraboots

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spilliams/terraboots/internal/generate"
	"github.com/spilliams/terraboots/internal/hclhelp"
)

type root struct {
	filename     string
	ID           string            `hcl:"id,label"`
	ScopeTypes   []string          `hcl:"scopeTypes"`
	Dependencies []*rootDependency `hcl:"dependency,block"`
	ScopeMatches []*scopeMatch     `hcl:"scopeMatch,block"`
}

type rootDependency struct {
	Root   string            `hcl:"root"`
	Scopes map[string]string `hcl:"scopes,optional"`
}

type scopeMatch struct {
	ScopeTypes map[string]string `hcl:"scopeTypes"`
}

func ParseRoot(cfgFile string) (*root, error) {
	// partial decode first, because we don't know what scope or attributes
	// this config will use. We're just looking for the `root` block here.
	cfg := &struct {
		Root *root `hcl:"root,block"`
	}{}

	err := hclsimple.DecodeFile(cfgFile, hclhelp.DefaultContext(), cfg)
	r := cfg.Root
	r.filename = cfgFile
	return r, err
}

func (p *Project) GenerateRoot(name string) error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.ScopeTypes))
	for i, el := range p.ScopeTypes {
		scopeTypes[i] = el.Name
	}

	return generate.Root(name, p.RootsDir, scopeTypes, p.Logger)
}
