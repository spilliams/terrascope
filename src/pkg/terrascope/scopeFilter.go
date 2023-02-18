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

// newScopeMatcher builds a new scopeFilterMatcher object
// address could be types & values interleaved, or just values
func newScopeMatcher(compiledScopes CompiledScopes, scopeTypes []*ScopeType, logger *logrus.Logger) *scopeMatcher {
	return &scopeMatcher{
		compiledScopes: compiledScopes,
		scopeTypes:     scopeTypes,

		Entry: logger.WithFields(logrus.Fields{"prefix": "scopeFilterMatcher"}),
	}
}

// determineMatchingScopes takes in a root configuration and an optional list of
// scopes. It returns a list of `CompiledScopes` where each scope in the
// list (a) matches at least one scopeMatch expression of the root, and
// (b) matches at least one scope given.
// Note that a root with no scopeMatch expressions will be treated as if all its
// scope types allow all values (`.*`).
func (sfm *scopeMatcher) determineMatchingScopes(root *root, scopes []string) (CompiledScopes, error) {
	matchingScopes := CompiledScopes{}
	// if they don't specify any scope matches, assume .* for all
	if root.ScopeMatches == nil || len(root.ScopeMatches) == 0 {
		allScopeMatchTypes := make(map[string]string)
		for _, scope := range root.ScopeTypes {
			allScopeMatchTypes[scope] = ".*"
		}
		root.ScopeMatches = []*scopeMatch{
			{ScopeTypes: allScopeMatchTypes},
		}
	}

	for _, scopeMatch := range root.ScopeMatches {
		matches, err := sfm.compiledScopes.Matching(scopeMatch.ScopeTypes)
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
			scopeFilter, err := sfm.makeFilter(scope)
			if err != nil {
				return nil, err
			}
			scopeFilters[i] = scopeFilter
		}
		sfm.Debugf("filters on the root's full list of scope values:\n%+v", scopeFilters)
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
			root.name, root.ScopeTypes, len(sfm.compiledScopes))
	}
	return matchingScopes, nil
}

func (sfm *scopeMatcher) makeFilter(address string) (map[string]string, error) {
	sfm.Debugf("makeScopeFilter %s", address)

	m := make(map[string]string)
	parts := strings.Split(address, ".")

	if len(parts)%2 == 0 {
		isCollated := true
		i := 0
		for i*2 < len(parts) {
			if parts[i*2] != sfm.scopeTypes[i].Name {
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
	sfm.Debugf("  parts after decollation: %v", parts)

	if len(parts) > len(sfm.scopeTypes) {
		return nil, fmt.Errorf("scope address %s is too long to be mapped to scope types %v", address, sfm.scopeTypes)
	}
	for i, v := range parts {
		// some special regex translating
		if v == "*" {
			v = ".*"
		}
		m[sfm.scopeTypes[i].Name] = v
	}
	sfm.Debugf("  mapping %+v", m)
	return m, nil
}
