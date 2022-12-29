package terraboots

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
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
