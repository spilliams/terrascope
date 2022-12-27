package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newScopeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scope",
		Aliases: []string{"s"},
		Short:   "Commands relating to scopes",
		GroupID: commandGroupIDTerraboots,

		PersistentPreRunE: bootsbootsPreRunE,
	}

	cmd.AddCommand(newScopeListCommand())
	cmd.AddCommand(newScopeGenerateCommand())

	return cmd
}

func newScopeListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "Lists all scope types in this project",
		RunE: func(cmd *cobra.Command, args []string) error {

			log.Infof("There are %d scopes in the project %s:", len(project.ScopeTypes), project.ID)
			for _, scope := range project.ScopeTypes {
				fmt.Println(scope.Name)
			}

			return nil
		},
	}

	return cmd
}

func newScopeGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g", "gen"},
		Short:   "Generates a new scope data file in this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.GenerateScopeData(os.Stdin, os.Stdout)
		},
	}

	return cmd
}
