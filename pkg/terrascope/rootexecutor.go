package terrascope

import (
	"path"
	"reflect"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
)

type RootDependencyChain int

const (
	RootDependencyChainUnknown RootDependencyChain = iota
	RootDependencyChainNone
	RootDependencyChainOne
	RootDependencyChainAll
)

type rootExecutorFactory struct {
	sm    *scopeMatcher
	rdc   *rootDependencyCalculator
	entry *logrus.Entry
}

func newRootExecutorFactory(sm *scopeMatcher, rdc *rootDependencyCalculator, logger *logrus.Logger) *rootExecutorFactory {
	return &rootExecutorFactory{
		sm:    sm,
		rdc:   rdc,
		entry: logger.WithFields(logrus.Fields{"prefix": "rootExec"}),
	}
}

// rootExecutor represents something that knows how to execute tasks on a root.
// This entails being able to enumerate all root-scope contexts applicable to
// the root as well as to handle how the root may depend on other roots.
type rootExecutor struct {
	root    *Root
	batches [][]*rootScopeContext
	*logrus.Entry
}

// newRootExecutor builds a new rootExecutor for the given root, scopes and
// other options.
// If `chain` is `RootExecutorDependencyChainingUnknown`, this function will
// survey the user for a "none/one/all" choice pertaining to the root's
// dependencies.
func (ref *rootExecutorFactory) newRootExecutor(root *Root, scopes []string, chain RootDependencyChain) (*rootExecutor, error) {
	// make sure we know how to handle dependencies (if we need to)
	if chain == RootDependencyChainUnknown && len(root.Dependencies) > 0 {
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
			chain = RootDependencyChainNone
		case one:
			chain = RootDependencyChainOne
		case all:
			chain = RootDependencyChainAll
		}
	}
	re := &rootExecutor{
		root:    root,
		batches: make([][]*rootScopeContext, 1),
		Entry:   ref.entry,
	}
	re.Debugf("root: %+v", root)

	batches, err := ref.rdc.prepareContextBatches(ref.sm, re.root, scopes, chain)
	if err != nil {
		return nil, err
	}

	re.batches = batches

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
		re.Infof("Note: This was a dry-run, so Terrascope can't guarantee to take exactly these actions if you re-run without the dry-run flag enabled.")
	}

	return outputs, nil
}
