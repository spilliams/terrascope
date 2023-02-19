package terrascope

import (
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

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

func newRootDependencyCalculator(roots map[string]*root, logger *logrus.Logger) *rootDependencyCalculator {
	return &rootDependencyCalculator{
		roots: roots,
		Entry: logger.WithField("prefix", "rootDepCalc"),
	}
}

type rootDependencyCalculator struct {
	roots map[string]*root
	*logrus.Entry
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

func (rdc *rootDependencyCalculator) prepareContextBatches(sm *scopeMatcher, r *root, scopes []string) ([][]*rootScopeContext, error) {
	// what scopes does the root apply to?
	matchingScopes, err := sm.determineMatchingScopes(r, scopes)
	if err != nil {
		return nil, err
	}
	rdc.Infof("Root will be executed for %d %s", len(matchingScopes), pluralize("scope", "scopes", len(matchingScopes)))
	for _, scope := range matchingScopes {
		rdc.Trace(scope.Address())
	}

	mainBatch := make([]*rootScopeContext, len(matchingScopes))
	for i, scope := range matchingScopes {
		mainBatch[i] = newRootScopeContext(r, scope, rdc.Logger)
	}
	// TODO root dependencies
	return nil, nil
}

// func (rdc *rootDependencyCalculator) prepareBatches(rootName string) ([][]string, error) {
// 	batchOrder := make(map[string]int)
// 	batchOrder[rootName] = 0
// 	chainAll := false
// 	switch rdc.chain {
// 	case RootDependencyChainNone:
// 		return makeBatches(batchOrder), nil
// 	case RootDependencyChainOne:
// 		break
// 	case RootDependencyChainAll:
// 		chainAll = true
// 		break
// 	default:
// 		return nil, fmt.Errorf("cannot prepare batches with unknown chaining rules")
// 	}

// 	return prepareBatchesRecursive(batchOrder, rootName, chainAll)
// }

// func (rdc *rootDependencyCalculator) prepareBatchesRecursive(batchOrder map[string]int, rootName string, recurse bool) ([][]string, error) {
// 	root, ok := rdc.roots[rootName]
// 	if !ok || root == nil {
// 		return nil, fmt.Errorf("cannot prepare batches for unrecognized root name '%s'", rootName)
// 	}
// 	rootOrder, ok := batchOrder[rootName]
// 	if !ok {
// 		return nil, fmt.Errorf("unknown error happened in the batch ordering process")
// 	}

// 	for {
// 		for _, dep := range root.Dependencies {

// 		}
// 		if !chainAll {
// 			return makeBatches(batchOrder, numBatches), nil
// 		}
// 	}
// }

func makeBatches(batchOrder map[string]int) [][]string {
	count := 0
	for _, index := range batchOrder {
		count = int(math.Max(float64(index), float64(count)))
	}
	batches := make([][]string, count)
	for i := 0; i < count; i++ {
		batches[i] = make([]string, 0)
	}
	for name, index := range batchOrder {
		batches[index] = append(batches[index], name)
	}
	return batches
}
