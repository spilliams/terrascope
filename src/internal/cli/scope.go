package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spilliams/terraboots/internal/scopedata"
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
	cmd.AddCommand(newScopeShowCommand())

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

func newScopeShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show SCOPE",
		Short: "Display a single scope value and it associated attributes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scopes, err := project.GetCompiledScopes(args[0])
			if err != nil {
				return err
			}
			if scopes == nil || len(scopes) == 0 {
				log.Errorf("No scope with address %s found", args[0])
				return nil
			}
			log.Infof("%d %s found with the scope filter '%s'", len(scopes), pluralize("scope", "scopes", len(scopes)), args[0])
			for _, scope := range scopes {
				printScope(scope)
			}
			return nil
		},
	}

	return cmd
}

func printScope(scope *scopedata.CompiledScope) {
	fmt.Println("Scope Details")
	for i := range scope.ScopeTypes {
		fmt.Printf("\t%s: %s\n", scope.ScopeTypes[i], scope.ScopeValues[i])
	}
	fmt.Println()
	fmt.Println("\tAttributes")
	for k, v := range scope.Attributes {
		fmt.Printf("\t%s: %s\n", k, v.AsString())
	}
}

func pluralize(single, plural string, count int) string {
	if count == 1 {
		return single
	}
	return plural
}
