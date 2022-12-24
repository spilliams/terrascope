package scopedata

// Scope represents a scope type
type Scope string

// Value represents one value for a single scope type. Instances may have
// children of the next scope type in the hierarchy.
type Value struct {
	Name     string
	Scope    Scope
	Address  string
	Children map[string]Value

	scopeTypeIndex int
}
