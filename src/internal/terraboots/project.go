package terraboots

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type Project struct {
	ID       string `hcl:"id,label"`
	RootsDir string `hcl:"rootsDir"`

	Scopes []*projectScope `hcl:"scope,block"`
	// roots  []*root
}

type projectScope struct {
	Name         string                    `hcl:"name"`
	Description  string                    `hcl:"description,optional"`
	DefaultValue string                    `hcl:"default,optional"`
	Validations  []*projectScopeValidation `hcl:"validation,block"`
}

type projectScopeValidation struct {
	Condition    bool   `hcl:"condition"`
	ErrorMessage string `hcl:"error_message"`
}

var terrabootsSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "terraboots", LabelNames: []string{"id"}},
	},
}

func ParseProject(filename string) (*Project, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}

	project := &Project{}
	content, diags := f.Body.Content(terrabootsSchema)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: map[string]function.Function{},
	}
	for _, block := range content.Blocks {
		if block.Type != "terraboots" {
			return nil, fmt.Errorf("found unknown block type %s: %+v", block.Type, block)
		}
		diags := gohcl.DecodeBody(block.Body, ctx, project)
		if err := handleDiags(diags, parser.Files(), nil); err != nil {
			return nil, err
		}
		project.ID = block.Labels[0]
	}

	return project, nil
}

func handleDiags(diags hcl.Diagnostics, files map[string]*hcl.File, writer io.Writer) error {
	if diags == nil {
		return nil
	}
	if writer == nil {
		writer = os.Stderr
	}
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			writer,
			files,
			100,   // wrapping width
			false, // colors
		)
		wr.WriteDiagnostics(diags)
		return fmt.Errorf("diagnostic errors found")
	}
	return nil
}
