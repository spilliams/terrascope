package ctyhelp

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

func String(v cty.Value) string {
	var vPrint string
	switch v.Type() {
	case cty.Number:
		if v.AsBigFloat().IsInt() {
			vPrint = fmt.Sprintf("%.0f", v.AsBigFloat())
		} else {
			vPrint = fmt.Sprintf("%f", v.AsBigFloat())
		}
	case cty.String:
		vPrint = v.AsString()
	case cty.Bool:
		vPrint = fmt.Sprintf("%v", v.True())
	default:
		vPrint = v.GoString()
	}
	return vPrint
}
