package terrascope

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zclconf/go-cty/cty"
)

// CompiledScope represents one permutation of several scope types. It contains
// attributes gained from each scope value along the way, with attribute values
// from narrower scopes overriding those of broader scopes
type CompiledScope struct {
	Attributes map[string]cty.Value

	// can't do map of types to values because maps are unordered
	ScopeTypes  []string
	ScopeValues []string
}

func (cs *CompiledScope) String() string {
	return cs.Address()
}

// Address returns the fully-qualified address of the receiver
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

// ToCtyValue returns a `cty.Value` representation of the receiver
func (cs *CompiledScope) ToCtyValue() cty.Value {
	kv := make(map[string]cty.Value)
	for i, scopeType := range cs.ScopeTypes {
		kv[scopeType] = cty.StringVal(cs.ScopeValues[i])
	}
	return cty.MapVal(kv)
}

// Values returns a period-separatec string of just the receiver's values.
func (cs *CompiledScope) Values() string {
	return strings.Join(cs.ScopeValues, ".")
}

// Matches returns true if the receiver matches the given types.
// The keys of `types` must match the scope types of the receiver.
// The values of `types` may be regular expressions.
// This function returns an error if and only if a value fails to
// `regexp.Compile`.
func (cs *CompiledScope) Matches(types map[string]string) (bool, error) {
	if len(types) != len(cs.ScopeTypes) {
		return false, nil
	}
	for matchKey, matchValue := range types {
		scopeIdx := indexOf(matchKey, cs.ScopeTypes)
		if scopeIdx == -1 {
			// different set of scope types were passed in
			return false, nil
		}
		myValue := cs.ScopeValues[scopeIdx]
		re, err := regexp.Compile(fmt.Sprintf("^%s$", matchValue))
		if err != nil {
			return false, err
		}
		if !re.MatchString(myValue) {
			return false, nil
		}
	}
	return true, nil
}

func indexOf(item string, list []string) int {
	for i, el := range list {
		if el == item {
			return i
		}
	}
	return -1
}

// CompiledScopes represents a list of CompiledScope objects.
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

// Deduplicate returns a CompiledScopes object that lists just the unique scopes
// of the receiver. Uniequeness is determined with the scopes' `Address` func.
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

// Matching returns a CompiledScopes object that lists the receiver's scopes
// that match the types given. Matches are determined with the scopes' `Matches`
// func.
func (css CompiledScopes) Matching(types map[string]string) (CompiledScopes, error) {
	newCSS := make([]*CompiledScope, 0, len(css))
	for _, scope := range css {
		ok, err := scope.Matches(types)
		if err != nil {
			return nil, err
		}
		if ok {
			newCSS = append(newCSS, scope)
		}
	}
	return newCSS, nil
}
