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

func newRootDependencyCalculator(roots map[string]*root, sm *scopeMatcher, logger *logrus.Logger) *rootDependencyCalculator {
	return &rootDependencyCalculator{
		roots: roots,
		sm:    sm,
		Entry: logger.WithField("prefix", "rootDepCalc"),
	}
}

type rootDependencyCalculator struct {
	roots map[string]*root
	sm    *scopeMatcher
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

func (rdc *rootDependencyCalculator) prepareContextBatches(sm *scopeMatcher, r *root, scopes []string, chainDependencies RootDependencyChain) ([][]*rootScopeContext, error) {
	if chainDependencies == RootDependencyChainUnknown {
		return nil, fmt.Errorf("cannot prepare batches with an unknown chaining rule")
	}
	rdc.Debugf("preparingContextBatches with chaining set to %d", chainDependencies)
	// map of rootName+scopeAddress to scope context object
	contexts := make(map[string]*rootScopeContext)
	// map of rootName+scopeAddress to order number (0-indexed, and reversed:
	// the higher the number, the earlier in the exec process it should go)
	batchOrder := make(map[string]int)
	// Contexts we have yet  to visit
	stack := make([]*rootScopeContext, 0)
	// first step is to prime the pump. base case.

	// what scopes does the root apply to?
	matchingScopes, err := sm.determineMatchingScopes(r, scopes)
	if err != nil {
		return nil, err
	}
	rdc.Infof("Root will be executed for %d %s", len(matchingScopes), pluralize("scope", "scopes", len(matchingScopes)))
	for _, scope := range matchingScopes {
		rdc.Trace(scope.Address())
		rootScope := newRootScopeContext(r, scope, rdc.Logger)
		contexts[rootScope.String()] = rootScope
		batchOrder[rootScope.String()] = 0
		stack = append(stack, rootScope)
	}
	highestOrder := 0

	// second step is to loop
	keepGoing := chainDependencies != RootDependencyChainNone
	for keepGoing && len(stack) > 0 {
		context := stack[0]
		key := context.String()
		order := batchOrder[key]

		for _, dep := range context.root.Dependencies {
			depRoot := rdc.roots[dep.RootName]
			// resolve the dep.RootName, dep.Scopes and context.scope into a new
			// compiled scope
			depScope, err := sm.resolveDependencyScope(depRoot, context.scope, dep.Scopes)
			if err != nil {
				return nil, err
			}
			depCtx := newRootScopeContext(depRoot, depScope, rdc.Logger)
			contexts[depCtx.String()] = depCtx
			// is it already in the order?
			if _, ok := batchOrder[depCtx.String()]; ok {
				newOrder := int(math.Max(float64(batchOrder[depCtx.String()]), float64(order+1)))
				batchOrder[depCtx.String()] = newOrder
				highestOrder = int(math.Max(float64(highestOrder), float64(newOrder)))
				continue
			}
			batchOrder[depCtx.String()] = order + 1
			highestOrder = int(math.Max(float64(highestOrder), float64(order+1)))
			// go look at these too if we should do more than just the direct deps
			if chainDependencies == RootDependencyChainAll {
				stack = append(stack, depCtx)
			}
		}
		stack = stack[1:]
	}

	// finally, turn the batch order map into a list (in the correct order)
	batchList := make([][]*rootScopeContext, highestOrder+1)
	for i := 0; i <= highestOrder; i++ {
		batchList[i] = make([]*rootScopeContext, 0)
	}
	for key, index := range batchOrder {
		reversedIndex := highestOrder - index
		batchList[reversedIndex] = append(batchList[reversedIndex], contexts[key])
	}

	for i, batch := range batchList {
		rdc.Infof("Batch %d:", i+1)
		for _, item := range batch {
			rdc.Infof("\t%s", item.String())
		}
	}
	return batchList, nil
}

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
