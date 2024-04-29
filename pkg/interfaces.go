package pkg

import "github.com/sirupsen/logrus"

type Logger interface {
	SetLogger(*logrus.Logger)
}

type Grapher interface {
	Logger
	RegisterReader(extension string, r GraphReader)
	RegisterWriter(extension string, w GraphWriter)
	ReadGraphFromFile(filename string) error
	WriteGraphToFile(filename string) error
	Graph() Graph
}

type GraphReader interface {
	Logger
	Read(filename string) (Graph, error)
}

type GraphWriter interface {
	Logger
	Write(g Graph, filename string) error
}

type Graph interface {
	Logger
	String() string
	Nodes() []Node
	// WalkNodes provides a way to iterate over the graph, operating on each node.
	// It takes in an iterator function that takes in a node of the graph, and
	// returns a node.
	// WalkNodes returns two integers. The first one should be the total number of
	// nodes walked (which should represent the number of nodes in the graph
	// *before* the walking), and the second one should be the sum of nodes
	// returned from the iterator (which should represent the number of nodes in
	// the graph *after* the walking).
	WalkNodes(func(Node) Node) (int, int)
	RemoveNode(string) error
	ChildToParents(name string) []string
}

type Node interface {
	String() string
	Attribute(AttributeKey) string
	SetAttribute(AttributeKey, string)
}

type AttributeKey interface {
	String() string
}
