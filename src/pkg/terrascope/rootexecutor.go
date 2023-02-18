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
	root    *root
	batches [][]*rootScopeContext
	*logrus.Entry

	// how to chain the dependencies found in the receiver's root contexts
	ChainDependencies RootExecutorDependencyChaining
}

// newRootExecutor builds a new rootExecutor for the given root, scopes and
// other options.
// If `chain` is `RootExecutorDependencyChainingUnknown`, this function will
// survey the user for a "none/one/all" choice pertaining to the root's
// dependencies.
func (p *Project) newRootExecutor(rootName string, scopes []string, chain RootExecutorDependencyChaining, logger *logrus.Logger) (*rootExecutor, error) {
	// make sure the root exists
	root, ok := p.Roots[rootName]
	if !ok {
		return nil, fmt.Errorf("Root '%s' isn't loaded. Did you run `AddAllRoots`?", rootName)
	}
	root = p.Roots[rootName]

	// make sure we know how to handle dependencies (if we need to)
	if chain == RootExecutorDependencyChainingUnknown && len(root.Dependencies) > 0 {
		var answer string
		none := "No, don't run any dependencies"
		one := "Yes, but just the direct dependencies"
		all := "Yes, and run all dependencies (direct and indirect)"
		err := survey.AskOne(&survey.Select{
			Message: "Run the same operation on the root's dependencies?",
			Options: []string{none, one, all},
			Default: none,
		}, &answer)
		if err != nil {
			return nil, err
		}
		switch answer {
		case none:
			chain = RootExecutorDependencyChainingNone
		case one:
			chain = RootExecutorDependencyChainingOne
		case all:
			chain = RootExecutorDependencyChainingAll
		}
	}
	re := &rootExecutor{
		root:    root,
		batches: make([][]*rootScopeContext, 1),
		Entry:   logger.WithFields(logrus.Fields{"prefix": "rootExec"}),

		ChainDependencies: chain,
	}
	re.Debugf("root: %+v", root)
	// set up the batches

	// what scopes does the root apply to?
	matchingScopes, err := p.determineMatchingScopes(root, scopes)
	if err != nil {
		return nil, err
	}
	re.Infof("Root will be executed for %d %s", len(matchingScopes), pluralize("scope", "scopes", len(matchingScopes)))
	for _, scope := range matchingScopes {
		re.Trace(scope.Address())
	}

	mainBatch := make([]*rootScopeContext, len(matchingScopes))
	for i, scope := range matchingScopes {
		mainBatch[i] = newRootScopeContext(root, scope, re.Entry.Logger)
	}
	re.batches[0] = mainBatch

	return re, nil
}

type ExecFunc func(*rootScopeContext) (string, error)

func (re *rootExecutor) Execute(f ExecFunc, dry bool) ([]string, error) {
	// TODO: use a worker pool
	// TODO: join all errors instead of exiting early
	outputs := make([]string, 0)

	for _, batch := range re.batches {
		for _, ctx := range batch {
			fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			fName = path.Base(fName)
			if dry {
				re.Debugf("would have run %s on %s", fName, ctx)
				continue
			} else {
				re.Debugf("running %s on %s", fName, ctx)
			}
			output, err := f(ctx)
			if err != nil {
				return nil, err
			}

			outputs = append(outputs, output)
		}
	}
	if dry {
		re.Infof("Note: This was a dry-run, so Terrascope can't guarantee to take exactly these actions if you re-run without the dry run flag enabled.")
	}

	return outputs, nil
}
