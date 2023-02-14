package graphviz

import (
	"io/ioutil"
	"os"

	"github.com/awalterschulze/gographviz"
	"github.com/sirupsen/logrus"
	"github.com/spilliams/tunnelvision/pkg"
)

type reader struct {
	*logrus.Logger
}

// NewReader returns a new file-reader that knows how to read graphviz
// files. It makes an assumption that the first node in the file is the root of
// the graph
func NewReader() pkg.GraphReader {
	return &reader{}
}

func (r *reader) SetLogger(l *logrus.Logger) {
	r.Logger = l
}

func (r *reader) Read(filename string) (pkg.Graph, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	graphAst, _ := gographviz.ParseString(string(contents))
	g := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, g); err != nil {
		return nil, err
	}
	graph := &graph{fundamental: g}
	graph.SetLogger(r.Logger)
	return graph, nil
}

type writer struct {
	*logrus.Logger
}

func NewWriter() pkg.GraphWriter {
	return &writer{}
}

func (w *writer) SetLogger(l *logrus.Logger) {
	w.Logger = l
}

func (w *writer) Write(g pkg.Graph, filename string) error {
	return os.WriteFile(filename, []byte(g.String()), 0777)
}
