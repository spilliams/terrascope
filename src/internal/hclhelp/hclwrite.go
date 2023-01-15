package hclhelp

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// CommentTokens helps write hcl tokens for a comment line in the configuration
func CommentTokens(msg string) hclwrite.Tokens {
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
