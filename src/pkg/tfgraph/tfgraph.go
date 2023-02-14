package tfgraph

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/tunnelvision/internal/graphviz"
	"github.com/spilliams/tunnelvision/pkg"
	"github.com/spilliams/tunnelvision/pkg/grapher"
)

func New(inFile string, logger *logrus.Logger, outFile string) error {
	gg := grapher.NewGrapher()
	gg.SetLogger(logger)
	gvReader := graphviz.NewReader()
	gg.RegisterReader("dot", gvReader)
	gg.RegisterReader("gv", gvReader)

	if err := gg.ReadGraphFromFile(inFile); err != nil {
		return err
	}

	// terraform makes graphs with some extra stuff. We should work on them one at
	// a time, so that later we could put them under feature flags.
	g := gg.Graph()
	_, after := g.WalkNodes(labelNode)
	logger.Infof("%d nodes labelled", after)
	before, after := g.WalkNodes(filterNode)
	logger.Infof("%d nodes filtered out", before-after)

	// TODO: best way to cluster?
	// before, after := g.WalkNodes(clusterNode)

	_, _ = g.WalkNodes(printNodeFunc(g, logger))

	gvWriter := graphviz.NewWriter()
	gg.RegisterWriter("dot", gvWriter)
	gg.RegisterWriter("gv", gvWriter)
	return gg.WriteGraphToFile(outFile)
}

func labelNode(n pkg.Node) pkg.Node {
	l := strings.TrimPrefix(n.String(), "\"")
	l = strings.TrimSuffix(l, "\"")
	l = strings.TrimPrefix(l, "[root] ")
	l = strings.TrimSuffix(l, " (expand)")
	l = strings.TrimSuffix(l, " (close)")
	n.SetAttribute(graphviz.LabelAttributeKey, l)
	return n
}

func filterNode(n pkg.Node) pkg.Node {
	switch typeOfNode(n) {
	case nodeTypeUnknown:
		logrus.Warnf("node type is unknown! dropping. node: %s", n)
		return nil
	// case nodeTypeVariable:
	// case nodeTypeOutput:
	// case nodeTypeModule:
	// case nodeTypeData:
	// case nodeTypeResource:
	case nodeTypeMeta:
		return nil
	// case nodeTypeProvider:
	// 	return nil
	case nodeTypeRoot:
		return nil
	default:
		return n
	}
}

func typeOfNode(n pkg.Node) nodeType {
	name := n.Attribute(graphviz.LabelAttributeKey)
	if strings.Contains(name, "provider[\\\"") {
		return nodeTypeProvider
	}
	if strings.HasPrefix(name, "meta.") {
		return nodeTypeMeta
	}
	if name == "root" {
		return nodeTypeRoot
	}

	parts := nameParts(name)
	if len(parts)%2 == 1 {
		return nodeTypeData
	}

	penult := parts[len(parts)-2]
	switch penult {
	case "var":
		return nodeTypeVariable
	case "output":
		return nodeTypeOutput
	case "local":
		return nodeTypeLocal
	case "module":
		return nodeTypeModule
	default:
		return nodeTypeResource
	}
}

type nodeType int

const (
	nodeTypeUnknown = iota
	nodeTypeVariable
	nodeTypeLocal
	nodeTypeOutput
	nodeTypeModule
	nodeTypeData
	nodeTypeResource
	nodeTypeMeta
	nodeTypeProvider
	nodeTypeRoot
)

func printNodeFunc(g pkg.Graph, logger *logrus.Logger) func(n pkg.Node) pkg.Node {
	return func(n pkg.Node) pkg.Node {
		logger.Debug(n)
		logger.Debugf("%v: %v", graphviz.LabelAttributeKey, n.Attribute(graphviz.LabelAttributeKey))
		logger.Debugf("parents: %v", g.ChildToParents(n.String()))
		logger.Debug()
		return n
	}
}

func nameParts(name string) []string {
	// TODO: save out the stuff that's inside brackets
	// split the leftover by "."
	// reinstate the brackets to each part
	return strings.Split(name, ".")
}

func clusterPath(name string) []string {
	// TODO
	return []string{}
}
