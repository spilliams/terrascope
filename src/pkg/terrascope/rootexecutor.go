package terrascope

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type RootExecutorDependencyChaining int

const (
	RootExecutorDependencyChainingUnknown RootExecutorDependencyChaining = iota
	RootExecutorDependencyChainingNone
	RootExecutorDependencyChainingOne
	RootExecutorDependencyChainingAll
)

type rootExecutor struct {
	root     *root
	contexts []*buildContext

	// how to chain the dependencies found in the receiver's root contexts
	ChainDependencies RootExecutorDependencyChaining
}

func (p *Project) newRootExecutor(rootName string, scopes []string) (*rootExecutor, error) {
	// make sure the root exists
	root, ok := p.Roots[rootName]
	if !ok {
		return nil, fmt.Errorf("Root '%s' isn't loaded. Did you run `AddAllRoots`?", rootName)
	}
	root = p.Roots[rootName]
	p.Debugf("root: %+v", root)
	re := &rootExecutor{root: root}

	// what scopes does the root apply to?
	matchingScopes, err := p.determineMatchingScopes(root, scopes)
	if err != nil {
		return nil, err
	}
	p.Infof("Root will be executed for %d %s", len(matchingScopes), pluralize("scope", "scopes", len(matchingScopes)))
	for _, scope := range matchingScopes {
		p.Trace(scope.Address())
	}

	builds := make([]*buildContext, len(matchingScopes))
	for i, scope := range matchingScopes {
		builds[i] = newBuildContext(root, scope, p.Entry.Logger)
	}
	re.contexts = builds

	if len(root.Dependencies) > 0 {
		var chain string
		none := "No, don't run any dependencies"
		one := "Yes, but just the direct dependencies"
		all := "Yes, and run all dependencies (direct and indirect)"
		err := survey.AskOne(&survey.Select{
			Message: "Run the same operation on the root's dependencies?",
			Options: []string{none, one, all},
			Default: none,
		}, &chain)
		if err != nil {
			return nil, err
		}
		switch chain {
		case none:
			re.ChainDependencies = RootExecutorDependencyChainingNone
		case one:
			re.ChainDependencies = RootExecutorDependencyChainingOne
		case all:
			re.ChainDependencies = RootExecutorDependencyChainingAll
		}
	}

	return re, nil
}

type ExecFunc func(*buildContext) (string, error)

func (re *rootExecutor) Execute(f ExecFunc) ([]string, error) {
	// TODO: use a worker pool
	outputs := make([]string, len(re.contexts))
	for i, build := range re.contexts {
		output, err := f(build)
		if err != nil {
			return nil, err
		}

		outputs[i] = output
	}

	return outputs, nil
}
