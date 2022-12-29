package scopedata

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// NestedScope represents one value for a single scope type. Instances may have
// children of the next scope type in the hierarchy.
// This is the type that we expect to read and write to "scope data" files.
type NestedScope struct {
	Type     string `hcl:"type,label"`
	Name     string `hcl:"name,label"`
	Address  string
	Children []*NestedScope `hcl:"scope,block"`
	Attrs    hcl.Attributes `hcl:",remain"`

	scopeTypeIndex int
}

// Count returns the number of NestedScopes known to the receiver, including
// itself.
func (ns *NestedScope) Count() int {
	if ns.Children == nil || len(ns.Children) == 0 {
		return 1
	}

	childCount := 0
	for _, child := range ns.Children {
		childCount += child.Count()
	}
	return 1 + childCount
}

// CompiledScope returns a CompiledScope object equivalent to the receiver by
// itself, without accounting for its children.
// Optionally provide a parent scope to
func (ns *NestedScope) CompiledScope(parent *CompiledScope) *CompiledScope {
	attrs := make(map[string]cty.Value)
	scopeTypes := make([]string, 0)
	scopeValues := make([]string, 0)
	if parent != nil {
		// TODO: scope value attributes; do a deep copy
		for k, v := range parent.Attributes {
			attrs[k] = v
		}
		scopeTypes = append(scopeTypes, parent.ScopeTypes...)
		scopeValues = append(scopeValues, parent.ScopeValues...)
	}

	scopeTypes = append(scopeTypes, ns.Type)
	scopeValues = append(scopeValues, ns.Name)

	for k, v := range ns.Attrs {
		// TODO: standard functions and variables?
		// what do with _ here?
		value, _ := v.Expr.Value(nil)
		attrs[k] = value
	}

	return &CompiledScope{
		Attributes:  attrs,
		ScopeTypes:  scopeTypes,
		ScopeValues: scopeValues,
	}
}

// CompiledScopes returns the complete set of CompiledScope objects equivalent
// to every permutation of the reciever and its children.
func (ns *NestedScope) CompiledScopes(parent *CompiledScope) []*CompiledScope {
	scopes := make([]*CompiledScope, 0, ns.Count())

	this := ns.CompiledScope(parent)
	scopes = append(scopes, this)
	for _, child := range ns.Children {
		scopes = append(scopes, child.CompiledScopes(this)...)
	}

	return scopes
}
