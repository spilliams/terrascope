package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/generate"
)

func newProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "Commands relating to a project",
		GroupID: commandGroupIDTerrascope,
	}

	cmd.AddCommand(newProjectGenerateCommand())
	cmd.AddCommand(newProjectGenerateScopesCommand())
	cmd.AddCommand(newProjectGraphRootDependenciesCommand())

	return cmd
}

func newProjectGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "g"},
		Short:   "Generates a new project in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate.Project(log.Logger)
		},
	}

	return cmd
}

func newProjectGenerateScopesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "generate-scopes",
		Short:             "Generates a new scope data file in this project",
		PersistentPreRunE: parseProject,
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.GenerateScopeData()
		},
	}

	return cmd
}

func newProjectGraphRootDependenciesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph-roots",
		Short: "Prints out a DOT-format graph of the roots in this Terrascope project and their dependencies",

		PersistentPreRunE: parseProject,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := project.AddAllRoots()
			if err != nil {
				return err
			}

			graph, err := project.RootDependencyGraph()
			if err != nil {
				return err
			}

			fmt.Println(graph)

			return nil
		},
	}
	return cmd
}
