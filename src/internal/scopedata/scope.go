package scopedata

// Scope represents one value for a single scope type. Instances may have
// children of the next scope type in the hierarchy.
type Scope struct {
	Type     string `hcl:"type,label"`
	Name     string `hcl:"name,label"`
	Address  string
	Children []*Scope `hcl:"scope,block"`

	scopeTypeIndex int
}

func (s *Scope) Count() int {
	if s.Children == nil || len(s.Children) == 0 {
		return 1
	}

	childCount := 0
	for _, child := range s.Children {
		childCount += child.Count()
	}
	return 1 + childCount
}
