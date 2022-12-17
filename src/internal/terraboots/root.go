package terraboots

import "github.com/hashicorp/hcl/v2/hclsimple"

type RootConfig struct {
	Root *Root `hcl:"root,block"`
}

type Root struct {
	ID           string            `hcl:"id,label"`
	Scopes       []RootScope       `hcl:"scopes"`
	Dependencies []*RootDependency `hcl:"dependency,block"`
}

type RootScope string

type RootDependency struct {
	Root   string               `hcl:"root"`
	Scopes map[RootScope]string `hcl:"scopes,optional"`
}

func ParseRoot(cfgFile string) (*Root, error) {
	cfg := &RootConfig{}
	err := hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Root, nil
}
