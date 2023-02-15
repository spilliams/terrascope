package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/terrascope/internal/hcl"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	rootDir := "/Users/spencer/spilliams/terrascope/src/fixtures/examples/simple-graph"
	logrus.Infof("reading configuration at %s", rootDir)
	// logrus.Infof("outputting graph in file %s", outFilename)

	parser := hcl.NewModule(logrus.StandardLogger())
	if err := parser.ParseModuleDirectory(rootDir); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// logrus.Debugf("parser: %#v", parser.Parser())
	// logrus.Debugf("module: %#v", parser.Module())
	// logrus.Debugf("configuration: %#v", parser.Configuration())

	graph, err := parser.DependencyGraph()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	graphJSON, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	fmt.Printf("dependencies: %v\n", string(graphJSON))

	logrus.Info("Done")
}
