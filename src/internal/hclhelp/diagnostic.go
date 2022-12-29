package hclhelp

import "github.com/hashicorp/hcl/v2"

func DiagnosticsWithoutSummary(err error, summary string) error {
	diags, typeOK := err.(hcl.Diagnostics)
	if !typeOK {
		return err
	}

	var newDiags hcl.Diagnostics
	for _, diag := range diags {
		if diag.Summary != summary {
			newDiags = append(newDiags, diag)
		}
	}

	if len(newDiags) > 0 {
		return newDiags
	}
	return nil
}
