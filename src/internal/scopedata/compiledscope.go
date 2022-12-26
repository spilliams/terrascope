package scopedata

import "fmt"

// CompiledScope represents one permutation of several scope types. It contains
// attributes gained from each scope value along the way, with attribute valuess
// from narrower scopes overriding those of broader scopes
type CompiledScope struct {
	Attributes map[string]interface{}

	// can't do map of types to values because maps are unordered
	ScopeTypes  []string
	ScopeValues []string
}

func (cs *CompiledScope) String() string {
	return cs.Address()
}

func (cs *CompiledScope) Address() string {
	addr := ""
	for i := range cs.ScopeTypes {
		addr += fmt.Sprintf("%s.%s", cs.ScopeTypes[i], cs.ScopeValues[i])
		if i < len(cs.ScopeTypes)-1 {
			addr += "."
		}
	}
	return addr
}

type CompiledScopes []*CompiledScope

func (css CompiledScopes) Len() int {
	return len(css)
}

func (css CompiledScopes) Less(i, j int) bool {
	return css[i].Address() < css[j].Address()
}

func (css CompiledScopes) Swap(i, j int) {
	css[i], css[j] = css[j], css[i]
}

func (css CompiledScopes) Deduplicate() []*CompiledScope {
	seen := make(map[string]bool)
	j := 0
	for _, scope := range css {
		if seen[scope.Address()] {
			continue
		}
		css[j] = scope
		j++
		seen[scope.Address()] = true
	}
	for i := j; i < len(css); i++ {
		css[i] = nil
	}
	css = css[:j]
	return css
}
