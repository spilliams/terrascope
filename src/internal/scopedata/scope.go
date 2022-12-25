package scopedata

// ScopeType represents a scope type
type ScopeType string

// Scope represents one value for a single scope type. Instances may have
// children of the next scope type in the hierarchy.
type Scope struct {
	Type     ScopeType `hcl:"type,label"`
	Name     string    `hcl:"name,label"`
	Address  string
	Children []*Scope

	scopeTypeIndex int
}
