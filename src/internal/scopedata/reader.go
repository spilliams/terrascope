package scopedata

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type Reader interface {
	Read() ([]Value, error)
}

type reader struct {
	files  []string
	scopes []string
	spec   hcldec.Spec
	*logrus.Logger
}

func NewReader(scopes []string, filenames []string, logger *logrus.Logger) Reader {
	return &reader{
		files:  filenames,
		scopes: scopes,
		Logger: logger,
	}
}

func (r *reader) Read() ([]Value, error) {
	if err := r.buildSpec(); err != nil {
		return nil, err
	}

	spec := r.spec.(*hcldec.BlockSpec)
	r.Debugf("reader spec: %+v", spec)
	for spec.Nested != nil {
		r.Debugf("   (cont'd): %+v", spec.Nested)
		spec = spec.Nested.(*hcldec.BlockSpec)
	}

	schema := hcldec.ImpliedSchema(r.spec)
	r.Debugf("schema: %+v", schema)

	values := make([]*cty.Value, len(r.files))
	for i, filename := range r.files {
		value, err := r.readScopeDataFile(filename)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return nil, nil
}

func (r *reader) buildSpec() error {
	spec := &hcldec.BlockSpec{}
	for i, scope := range r.scopes {
		spec.TypeName = scope

		if r.spec == nil {
			spec.Required = true
			r.spec = spec
		}

		if i < len(r.scopes)-1 {
			parent := spec
			child := &hcldec.BlockSpec{}
			parent.Nested = child
			spec = child
		}
	}
	return nil
}

// readScopeDataFile reads a single file containing scope data
func (r *reader) readScopeDataFile(filename string) (*cty.Value, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}

	// content, diags := f.Body.Content(schema)
	// if err := handleDiags(diags, parser.Files(), nil); err != nil {
	// 	return err
	// }
	// ctx := &hcl.EvalContext{
	// 	Variables: map[string]cty.Value{},
	// 	Functions: map[string]function.Function{},
	// }
	value, diags := hcldec.Decode(f.Body, r.spec, nil)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}

	return &value, nil
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
