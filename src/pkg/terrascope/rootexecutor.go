package terrascope

import (
	"fmt"
	"path"
	"reflect"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
)

type RootExecutorDependencyChaining int

const (
	RootExecutorDependencyChainingUnknown RootExecutorDependencyChaining = iota
	RootExecutorDependencyChainingNone
	RootExecutorDependencyChainingOne
	RootExecutorDependencyChainingAll
)

// rootExecutor represents something that knows how to execute tasks on a root.
// This entails being able to enumerate all root-scope contexts applicable to
// the root as well as to handle how the root may depend on other roots.
type rootExecutor struct {
	root     *root
	contexts []*rootScopeContext
	*logrus.Entry

	// how to chain the dependencies found in the receiver's root contexts
	ChainDependencies RootExecutorDependencyChaining
}

func (p *Project) newRootExecutor(rootName string, scopes []string, logger *logrus.Logger) (*rootExecutor, error) {
	// make sure the root exists
	root, ok := p.Roots[rootName]
	if !ok {
		return nil, fmt.Errorf("Root '%s' isn't loaded. Did you run `AddAllRoots`?", rootName)
	}
	root = p.Roots[rootName]
	p.Debugf("root: %+v", root)
	re := &rootExecutor{
		root:  root,
		Entry: logger.WithFields(logrus.Fields{"prefix": "rootExec"}),
	}

	// what scopes does the root apply to?
	matchingScopes, err := p.determineMatchingScopes(root, scopes)
	if err != nil {
		return nil, err
	}
	p.Infof("Root will be executed for %d %s", len(matchingScopes), pluralize("scope", "scopes", len(matchingScopes)))
	for _, scope := range matchingScopes {
		p.Trace(scope.Address())
	}

	builds := make([]*rootScopeContext, len(matchingScopes))
	for i, scope := range matchingScopes {
		builds[i] = newRootScopeContext(root, scope, p.Entry.Logger)
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

type ExecFunc func(*rootScopeContext) (string, error)

func (re *rootExecutor) Execute(f ExecFunc, dry bool) ([]string, error) {
	// TODO: use a worker pool
	// TODO: join all errors instead of exiting early
	outputs := make([]string, 0)
	for _, build := range re.contexts {
		fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		fName = path.Base(fName)
		if dry {
			re.Debugf("would have run %s on %s", fName, build)
			continue
		} else {
			re.Debugf("running %s on %s", fName, build)
		}
		output, err := f(build)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, output)
	}
	if dry {
		re.Infof("Note: This was a dry-run, so Terrascope can't guarantee to take exactly these actions if you re-run without the dry run flag enabled.")
	}

	return outputs, nil
}
