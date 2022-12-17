package terraboots

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type ProjectConfig struct {
	Project *Project `hcl:"terraboots,block"`
}

type Project struct {
	ID       string `hcl:"id,label"`
	RootsDir string `hcl:"rootsDir"`

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
	// TODO: use this for parsing scope validation blocks...
	// ctx := &hcl.EvalContext{
	// 	Variables: map[string]cty.Value{
	// 		"scope": cty.ObjectVal(map[string]cty.Value{
	// 			"value":
	// 		})
	// 	},
	// 	// shamelessly stolen from Terraform
	// 	// https://github.com/hashicorp/terraform/blob/f8669d235174ebdaf503d1cd400e22eb51c74c3b/internal/lang/functions.go
	// 	Functions: map[string]function.Function{
	// 		"abs":             stdlib.AbsoluteFunc,
	// 		"can":             tryfunc.CanFunc,
	// 		"ceil":            stdlib.CeilFunc,
	// 		"chomp":           stdlib.ChompFunc,
	// 		"coalescelist":    stdlib.CoalesceListFunc,
	// 		"compact":         stdlib.CompactFunc,
	// 		"concat":          stdlib.ConcatFunc,
	// 		"contains":        stdlib.ContainsFunc,
	// 		"csvdecode":       stdlib.CSVDecodeFunc,
	// 		"distinct":        stdlib.DistinctFunc,
	// 		"element":         stdlib.ElementFunc,
	// 		"chunklist":       stdlib.ChunklistFunc,
	// 		"flatten":         stdlib.FlattenFunc,
	// 		"floor":           stdlib.FloorFunc,
	// 		"format":          stdlib.FormatFunc,
	// 		"formatdate":      stdlib.FormatDateFunc,
	// 		"formatlist":      stdlib.FormatListFunc,
	// 		"indent":          stdlib.IndentFunc,
	// 		"join":            stdlib.JoinFunc,
	// 		"jsondecode":      stdlib.JSONDecodeFunc,
	// 		"jsonencode":      stdlib.JSONEncodeFunc,
	// 		"keys":            stdlib.KeysFunc,
	// 		"log":             stdlib.LogFunc,
	// 		"lower":           stdlib.LowerFunc,
	// 		"max":             stdlib.MaxFunc,
	// 		"merge":           stdlib.MergeFunc,
	// 		"min":             stdlib.MinFunc,
	// 		"parseint":        stdlib.ParseIntFunc,
	// 		"pow":             stdlib.PowFunc,
	// 		"range":           stdlib.RangeFunc,
	// 		"regex":           stdlib.RegexFunc,
	// 		"regexall":        stdlib.RegexAllFunc,
	// 		"reverse":         stdlib.ReverseListFunc,
	// 		"setintersection": stdlib.SetIntersectionFunc,
	// 		"setproduct":      stdlib.SetProductFunc,
	// 		"setsubtract":     stdlib.SetSubtractFunc,
	// 		"setunion":        stdlib.SetUnionFunc,
	// 		"signum":          stdlib.SignumFunc,
	// 		"slice":           stdlib.SliceFunc,
	// 		"sort":            stdlib.SortFunc,
	// 		"split":           stdlib.SplitFunc,
	// 		"strrev":          stdlib.ReverseFunc,
	// 		"substr":          stdlib.SubstrFunc,
	// 		"timeadd":         stdlib.TimeAddFunc,
	// 		"title":           stdlib.TitleFunc,
	// 		"trim":            stdlib.TrimFunc,
	// 		"trimprefix":      stdlib.TrimPrefixFunc,
	// 		"trimspace":       stdlib.TrimSpaceFunc,
	// 		"trimsuffix":      stdlib.TrimSuffixFunc,
	// 		"try":             tryfunc.TryFunc,
	// 		"upper":           stdlib.UpperFunc,
	// 		"values":          stdlib.ValuesFunc,
	// 		"zipmap":          stdlib.ZipmapFunc,
	// 	},
	// }
	err := hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}

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
