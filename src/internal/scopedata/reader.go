package scopedata

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sirupsen/logrus"
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

	spec := r.spec.(*hcldec.BlockMapSpec)
	r.Debugf("reader spec: %+v", spec)
	for spec.Nested != nil {
		spec = spec.Nested.(*hcldec.BlockMapSpec)
		r.Debugf("   (cont'd): %+v", spec)
	}

	schema := hcldec.ImpliedSchema(r.spec)
	r.Debugf("schema: %#v", schema)

	values := make([]*Value, 0)
	for _, filename := range r.files {
		fileValues, err := r.readScopeDataFile(filename)
		if err != nil {
			return nil, err
		}
		values = append(values, fileValues...)
	}

	return nil, nil
}

func (r *reader) buildSpec() error {
	spec := &hcldec.BlockMapSpec{}
	for i, scope := range r.scopes {
		spec.TypeName = scope
		spec.LabelNames = []string{"id"}

		if r.spec == nil {
			r.spec = spec
		}

		if i < len(r.scopes)-1 {
			child := &hcldec.BlockMapSpec{}
			spec.Nested = child
			spec = child
		}
	}
	return nil
}

// readScopeDataFile reads a single file containing scope data
func (r *reader) readScopeDataFile(filename string) ([]*Value, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}
	r.Debugf("scope data file body: %+v", f.Body)

	schema := hcldec.ImpliedSchema(r.spec)
	content, partial, diags := f.Body.PartialContent(schema)
	// content, diags := f.Body.Content(schema)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}
	r.Debugf("f body content: %+v", content)
	r.Debugf("f body partial: %+v", partial)
	// TODO: the real data is still in the partial. The implied schema did not
	// pull anything out of the file.

	_, diags = hcldec.Decode(f.Body, r.spec, nil)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}

	// TODO turn cty values into scope values...

	return nil, nil
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
