package terraboots

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

type rootConfig struct {
	Root       *root        `hcl:"root,block"`
	Generators []*generator `hcl:"generate,block"`
	Includes   []*include   `hcl:"include,block"`
	Inputs     *cty.Value   `hcl:"inputs,attr"`
}

type root struct {
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

type generator struct {
	ID       string `hcl:"id,label"`
	Path     string `hcl:"path,attr"`
	Contents string `hcl:"contents,attr"`
}

type include struct {
	Path string `hcl:"path,attr"`
}

func ParseRoot(cfgFile string) (*root, error) {
	cfg := &rootConfig{}
	// TODO: build a root; we need a more advanced decode here, to allow for
	// partial decoding and also Functions & Values
	err := hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Root, nil
}
