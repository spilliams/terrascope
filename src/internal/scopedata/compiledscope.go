package scopedata

import (
	"fmt"
	"regexp"
	"strings"
)

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

func (cs *CompiledScope) Values() string {
	return strings.Join(cs.ScopeValues, ".")
}

func (cs *CompiledScope) Matches(types map[string]string) bool {
	if len(types) != len(cs.ScopeTypes) {
		return false
	}
	for matchKey, matchValue := range types {
		scopeIdx := indexOf(matchKey, cs.ScopeTypes)
		if scopeIdx == -1 {
			// different set of scope types were passed in
			return false
		}
		myValue := cs.ScopeValues[scopeIdx]
		re := regexp.MustCompile(matchValue)
		if !re.MatchString(myValue) {
			return false
		}
	}
	return true
}

func indexOf(item string, list []string) int {
	for i, el := range list {
		if el == item {
			return i
		}
	}
	return -1
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

func (css CompiledScopes) Deduplicate() CompiledScopes {
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

func (css CompiledScopes) Matching(types map[string]string) CompiledScopes {
	newCSS := make([]*CompiledScope, 0, len(css))
	for _, scope := range css {
		if scope.Matches(types) {
			newCSS = append(newCSS, scope)
		}
	}
	return newCSS
}
