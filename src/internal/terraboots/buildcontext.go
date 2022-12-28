package terraboots

import (
	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/internal/scopedata"
)

type buildContext struct {
	root  *root
	scope *scopedata.CompiledScope
	*logrus.Entry
}

func newBuildContext(root *root, scope *scopedata.CompiledScope, logger *logrus.Logger) *buildContext {
	return &buildContext{
		root:  root,
		scope: scope,
		Entry: logger.WithFields(logrus.Fields{
			"prefix": "builder",
			"root":   root.ID,
			"scope":  scope.Address(),
		}),
	}
}

func (bc *buildContext) Build() {
	bc.Info("Building root")
}
