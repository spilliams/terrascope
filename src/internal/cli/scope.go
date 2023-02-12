package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/ctyhelp"
	"github.com/spilliams/terrascope/pkg/terrascope"
)

func newScopeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scope",
		Aliases: []string{"s"},
		Short:   "Commands relating to scopes",
		GroupID: commandGroupIDTerrascope,

		PersistentPreRunE: parseProject,
	}

	cmd.AddCommand(newScopeListCommand())
	// cmd.AddCommand(newScopeGenerateCommand())
	cmd.AddCommand(newScopeShowCommand())

	return cmd
}

func newScopeListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "Lists all scope types in this project",
		RunE: func(cmd *cobra.Command, args []string) error {

			log.Infof("There %s %d scope %s in the project %s:",
				pluralize("is", "are", len(project.ScopeTypes)),
				len(project.ScopeTypes),
				pluralize("type", "types", len(project.ScopeTypes)),
				project.ID)
			for _, scope := range project.ScopeTypes {
				fmt.Println(scope.Name)
			}

			return nil
		},
	}

	return cmd
}

// func newScopeGenerateCommand() *cobra.Command {
// 	(reserved for generating a single new scope)
// }

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

func printScope(scope *terrascope.CompiledScope) {
	fmt.Printf("%s:\n", scope.Address())
	for k, v := range scope.Attributes {
		vPrint := ctyhelp.String(v)
		fmt.Printf("\t%s: %s\n", k, vPrint)
	}
}

func pluralize(single, plural string, count int) string {
	if count == 1 {
		return single
	}
	return plural
}
