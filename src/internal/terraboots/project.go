package terraboots

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type ProjectConfig struct {
	Project *Project `hcl:"terraboots,block"`
}

type Project struct {
	ID       string `hcl:"id,label"`
	RootsDir string `hcl:"rootsDir"`

	Scopes []*ProjectScope `hcl:"scope,block"`
	// roots  []*root
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

func ParseProject(filename string) (*Project, error) {
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
	err := hclsimple.DecodeFile(filename, nil, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Project, nil
}
