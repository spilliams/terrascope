package hcl

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
)

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
