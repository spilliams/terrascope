package hcl

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// TokensForComment helps write hcl tokens for a comment line in the configuration
func TokensForComment(msg string) hclwrite.Tokens {
	if !strings.HasPrefix(msg, "# ") {
		msg = fmt.Sprintf("# %s", msg)
	}
	msgToken := &hclwrite.Token{
		Type:         hclsyntax.TokenComment,
		Bytes:        []byte(msg),
		SpacesBefore: 0,
	}
	return []*hclwrite.Token{msgToken}
}

// TokensForObjectWithComments returns Tokens representing an hcl object of
// the given attributes, with the given comments appended to the end of each
// attribute line
func TokensForObjectWithComments(attrs map[string]cty.Value, comments map[string]string) hclwrite.Tokens {
	var toks hclwrite.Tokens
	toks = append(toks, &hclwrite.Token{
		Type:  hclsyntax.TokenOBrace,
		Bytes: []byte{'{'},
	})
	toks = append(toks, &hclwrite.Token{
		Type:  hclsyntax.TokenNewline,
		Bytes: []byte{'\n'},
	})
	for key, val := range attrs {
		toks = append(toks, hclwrite.TokensForValue(cty.StringVal(key))...)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenEqual,
			Bytes: []byte{'='},
		})
		toks = append(toks, hclwrite.TokensForValue(val)...)

		if comment, ok := comments[key]; ok {
			toks = append(toks, &hclwrite.Token{
				Type:         hclsyntax.TokenComment,
				Bytes:        []byte(fmt.Sprintf("# from %s", comment)),
				SpacesBefore: 0,
			})
		}

		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		})
	}
	toks = append(toks, &hclwrite.Token{
		Type:  hclsyntax.TokenCBrace,
		Bytes: []byte{'}'},
	})

	return toks
}
