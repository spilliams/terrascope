package graphviz

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/spilliams/tunnelvision/pkg"
)

type node struct {
	fundamental *gographviz.Node
}

func (n *node) String() string {
	return n.fundamental.Name
}

func (n *node) Attribute(key pkg.AttributeKey) string {
	val, ok := n.fundamental.Attrs[gographviz.Attr(key.String())]
	if !ok {
		return val
	}
	return strings.TrimPrefix(strings.TrimSuffix(val, `"`), `"`)
}

func (n *node) SetAttribute(key pkg.AttributeKey, value string) {
	strings.Join(strings.Split(value, `"`), `\"`)
	n.fundamental.Attrs.Add(key.String(), fmt.Sprintf(`"%s"`, value))
}
