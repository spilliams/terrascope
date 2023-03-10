package terrascope

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type scopeMatcher struct {
	compiledScopes CompiledScopes
	scopeTypes     []*ScopeType
	*logrus.Entry
}

// newScopeMatcher builds a new scopeMatcher object
// address could be types & values interleaved, or just values
func newScopeMatcher(compiledScopes CompiledScopes, scopeTypes []*ScopeType, logger *logrus.Logger) *scopeMatcher {
	return &scopeMatcher{
		compiledScopes: compiledScopes,
		scopeTypes:     scopeTypes,

		Entry: logger.WithFields(logrus.Fields{"prefix": "scopeMatcher"}),
	}
}

// determineMatchingScopes takes in a root configuration and an optional list of
// scopes. It returns a list of `CompiledScopes` where each scope in the
// list (a) matches at least one scopeMatch expression of the root, and
// (b) matches at least one scope given.
// Note that a root with no scopeMatch expressions will be treated as if all its
// scope types allow all values (`.*`).
func (sm *scopeMatcher) determineMatchingScopes(root *Root, scopes []string) (CompiledScopes, error) {
	matchingScopes := CompiledScopes{}
	// if they don't specify any scope matches, assume .* for all
	if root.ScopeMatches == nil || len(root.ScopeMatches) == 0 {
		allScopeMatchTypes := make(map[string]string)
		for _, scope := range root.ScopeTypes {
			allScopeMatchTypes[scope] = ".*"
		}
		root.ScopeMatches = []*ScopeMatch{
			{ScopeTypes: allScopeMatchTypes},
		}
	}

	for _, scopeMatch := range root.ScopeMatches {
		matches, err := sm.compiledScopes.Matching(scopeMatch.ScopeTypes)
		if err != nil {
			return nil, err
		}
		matchingScopes = append(matchingScopes, matches...)
	}
	matchingScopes = matchingScopes.Deduplicate()
	sort.Sort(matchingScopes)

	// also abide by this list
	if len(scopes) > 0 {
		filteredMatchingScopes := CompiledScopes{}
		scopeFilters := make([]map[string]string, len(scopes))
		for i, scope := range scopes {
			scopeFilter, err := sm.makeFilter(scope)
			if err != nil {
				return nil, err
			}
			scopeFilters[i] = scopeFilter
		}
		sm.Debugf("filters on the root's full list of scope values:\n%+v", scopeFilters)
		for _, scope := range matchingScopes {
			for _, filter := range scopeFilters {
				ok, err := scope.Matches(filter)
				if err != nil {
					return nil, err
				}
				if ok {
					filteredMatchingScopes = append(filteredMatchingScopes, scope)
				}
			}
		}
		matchingScopes = filteredMatchingScopes
	}

	if len(matchingScopes) == 0 {
		return nil, fmt.Errorf("No matching scope values found.\n"+
			"Root '%s' applies to the scope types %v.\n"+
			"All scopes in the project were searched (%d), and none matched these types. Please provide\n"+
			"\t- new scope data for the project,\n"+
			"\t- different scope types in the root configuration file, or\n"+
			"\t- new scope matches in the root configuration file.",
			root.Name, root.ScopeTypes, len(sm.compiledScopes))
	}
	return matchingScopes, nil
}

func (sm *scopeMatcher) makeFilter(address string) (map[string]string, error) {
	sm.Debugf("makeScopeFilter %s", address)

	m := make(map[string]string)
	parts := strings.Split(address, ".")

	if len(parts)%2 == 0 {
		isCollated := true
		i := 0
		for i*2 < len(parts) {
			if parts[i*2] != sm.scopeTypes[i].Name {
				isCollated = false
				break
			}
			i++
		}

		if isCollated {
			newParts := make([]string, 0)
			for i := 1; i < len(parts); i += 2 {
				newParts = append(newParts, parts[i])
			}
			parts = newParts
		}
	}
	sm.Debugf("  parts after decollation: %v", parts)

	if len(parts) > len(sm.scopeTypes) {
		return nil, fmt.Errorf("scope address %s is too long to be mapped to scope types %v", address, sm.scopeTypes)
	}
	for i, v := range parts {
		// some special regex translating
		if v == "*" {
			v = ".*"
		}
		m[sm.scopeTypes[i].Name] = v
	}
	sm.Debugf("  mapping %+v", m)
	return m, nil
}

func (sm *scopeMatcher) resolveDependencyScope(ancestor *Root, descendantScope *CompiledScope, ancestorScopes map[string]string) (*CompiledScope, error) {
	ancestorFullScopes := make(map[string]string)
	for _, scopeType := range ancestor.ScopeTypes {
		if customValue, ok := ancestorScopes[scopeType]; ok {
			ancestorFullScopes[scopeType] = customValue
			continue
		}
		descendantScopeTypeIndex := indexOf(scopeType, descendantScope.ScopeTypes)
		descendantScopeValue := descendantScope.ScopeValues[descendantScopeTypeIndex]
		ancestorFullScopes[scopeType] = descendantScopeValue
	}
	sm.Tracef("Searching compiled scopes for a match to the scope description %v", ancestorFullScopes)
	matches := make([]*CompiledScope, 0)
	for _, cs := range sm.compiledScopes {
		match, err := cs.Matches(ancestorFullScopes)
		if err != nil {
			return nil, err
		}
		if match {
			sm.Tracef("\tmatch found: %s", cs.String())
			matches = append(matches, cs)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("cannot find a compiled scope matching the description %v", ancestorFullScopes)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple matches found for a single scope description %v", ancestorFullScopes)
	}
	return matches[0], nil
}
