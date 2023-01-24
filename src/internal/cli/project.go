package cli

import (
	"github.com/spf13/cobra"
	"github.com/spilliams/terraboots/internal/generate"
)

func newProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "Commands relating to a project",
		GroupID: commandGroupIDTerraboots,
	}

	cmd.AddCommand(newProjectGenerateCommand())
	cmd.AddCommand(newProjectGenerateScopesCommand())

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
		PersistentPreRunE: bootsbootsPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.GenerateScopeData()
		},
	}

	return cmd
}
