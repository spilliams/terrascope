package terraboots

import "github.com/hashicorp/hcl/v2/hclsimple"

type RootConfig struct {
	Root *Root `hcl:"root,block"`
}

type Root struct {
	ID           string            `hcl:"id,label"`
	ScopeTypes   []string          `hcl:"scopeTypes"`
	Dependencies []*RootDependency `hcl:"dependency,block"`
	ScopeMatches []*ScopeMatch     `hcl:"scopeMatch,block"`
}

type RootDependency struct {
	Root   string            `hcl:"root"`
	Scopes map[string]string `hcl:"scopes,optional"`
}

type ScopeMatch struct {
	ScopeTypes map[string]string `hcl:"scopeTypes"`
}

func ParseRoot(cfgFile string) (*Root, error) {
	cfg := &RootConfig{}
	err := hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Root, nil
}
