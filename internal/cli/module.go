package cli

import (
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
	}

	cmd.AddCommand(newModuleGraphResourcesCommand())

	return cmd
}

func newModuleGraphResourcesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "graph-resources [DIRECTORY]",
		Short: "(EXPERIMENTAL) graphs the root module at the given directory (`.` by default)",
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
			return printModuleGraph(rootDir)
		},
	}
}

func printModuleGraph(dir string) error {
	log.Infof("reading configuration at %s", dir)

	parser := hcl.NewModule(log.Logger)
	if err := parser.ParseModuleDirectory(dir); err != nil {
		return err
	}

	graph, err := parser.DependencyGraph()
	if err != nil {
		return err
	}
	fmt.Println(graph)
	log.Warnf("Note: this graph is experimental and may not represent\n100%% of the resources or relationships of your configuration. If you know how\nthis could be improved, please submit an Issue or a PR to the source repository!\nhttps://github.com/spilliams/terrascope")
	return nil
}
