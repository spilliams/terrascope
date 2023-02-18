package terrascope

import "fmt"

type root struct {
	filename     string
	name         string
	ScopeTypes   []string          `hcl:"scopeTypes"`
	Dependencies []*rootDependency `hcl:"dependency,block"`
	ScopeMatches []*scopeMatch     `hcl:"scopeMatch,block"`
}

type rootDependency struct {
	RootName string            `hcl:"root"`
	Scopes   map[string]string `hcl:"scopes,optional"`
}

type scopeMatch struct {
	ScopeTypes map[string]string `hcl:"scopeTypes"`
}

type rootDependencyCalculator struct {
	roots map[string]*root
	chain RootExecutorDependencyChaining
}

func (rdc *rootDependencyCalculator) assertRootDependenciesAcyclic() error {
	visited := make(map[string]bool)
	for rootName := range rdc.roots {
		if visited[rootName] {
			continue
		}
		stack := make(map[string]bool)
		if has, list := rdc.hasCyclicDependency(rootName, visited, stack); has {
			return fmt.Errorf("cyclical dependency detected for root '%s': %+v", rootName, list)
		}
	}
	return nil
}

// hasCyclicDependency performs a depth-first search. It takes the current root
// name, the map of all roots, a visited map to keep track of which nodes have
// been visited, and a stack map to keep track of nodes we have yet to visit.
// It returns true if a cyclical dependency is found, and false otherwise.
// When it returns true, it also returns the list of root names in the cycle.
func (rdc *rootDependencyCalculator) hasCyclicDependency(rootName string, visited, stack map[string]bool) (bool, []string) {
	visited[rootName] = true
	stack[rootName] = true

	for _, dep := range rdc.roots[rootName].Dependencies {
		if !visited[dep.RootName] {
			if has, _ := rdc.hasCyclicDependency(dep.RootName, visited, stack); has {
				return true, keys(visited)
			}
		} else if stack[dep.RootName] {
			return true, keys(visited)
		}
	}

	delete(stack, rootName)
	return false, []string{}
}

func keys(m map[string]bool) []string {
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}
	return l
}

func (rdc *rootDependencyCalculator) prepareBatches(rootName string) error {
	batchOrder := make(map[string]int)
	if err := rdc.prepareBatch(batchOrder, 0); err != nil {
		return err
	}
	batchOrder[rootName] = 0
	switch rdc.chain {
	case RootExecutorDependencyChainingNone:
	case RootExecutorDependencyChainingOne:
	case RootExecutorDependencyChainingAll:
	}
	// TODO root dependencies
	return nil
}

func (rdc *rootDependencyCalculator) prepareBatch(batchOrder map[string]int, cursor int) error {
	// TODO root dependencies
	return nil
}
