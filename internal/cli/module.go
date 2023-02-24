package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/hcl"
)

func newModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "module",
		Aliases: []string{"m"},
		Short:   "A toolbox for working with Terraform modules",
		GroupID: commandGroupIDTunnelvision,
	}

	cmd.AddCommand(newModuleGraphResourcesCommand())

	return cmd
}

func newModuleGraphResourcesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "graph-resources [DIRECTORY]",
		Short: "graphs the root module at the given directory (`.` by default)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir, err := os.Getwd()
			if err != nil {
				return err
			}
			if len(args) > 0 {
				rootDir = args[0]
				rootDir, err = filepath.Abs(rootDir)
				if err != nil {
					return err
				}
			}
			// outFilename := "output.dot"
			graph, err := graphModule(rootDir)
			if err != nil {
				return err
			}
			fmt.Println(string(graph))
			return nil
		},
	}
}

func graphModule(dir string) ([]byte, error) {
	log.Infof("reading configuration at %s", dir)
	// logrus.Infof("outputting graph in file %s", outFilename)

	parser := hcl.NewModule(log.Logger)
	if err := parser.ParseModuleDirectory(dir); err != nil {
		return nil, err
	}

	// logrus.Debugf("%#v", parser.Parser())
	// logrus.Debugf("%#v", parser.Module())
	// logrus.Debugf("configuration: %#v", parser.Configuration())

	// TODO: build a graph from the module

	graph, err := parser.DependencyGraph()
	if err != nil {
		return nil, err
	}

	graphJSON, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return nil, err
	}
	return graphJSON, err
}
