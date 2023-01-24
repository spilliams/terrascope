package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "root",
		Aliases: []string{"r"},
		Short:   "Commands relating to root modules",
		GroupID: commandGroupIDTerraboots,

		PersistentPreRunE: bootsbootsPreRunE,
	}

	cmd.AddCommand(newRootBuildCommand())
	cmd.AddCommand(newRootGenerateCommand())
	// cmd.AddCommand(newRootGraphCommand())
	cmd.AddCommand(newRootListCommand())

	return cmd
}

func newRootBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build ROOT [SCOPE]...",
		Aliases: []string{"b"},
		Short:   "Builds the given root and prints the location of the built configurations to stdout",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scopes := make([]string, len(args)-1)
			for i := 1; i < len(args); i++ {
				scopes[i-1] = args[i]
			}
			dirs, err := project.BuildRoot(args[0], scopes)

			for _, dir := range dirs {
				fmt.Println(dir)
			}
			return err
		},
	}

	return cmd
}

func newRootGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate [NAME]",
		Aliases: []string{"gen", "g"},
		Short:   "Generates a new root module",
		RunE: func(cmd *cobra.Command, args []string) error {
			var rootName string
			if len(args) > 0 {
				rootName = args[0]
			}
			return project.GenerateRoot(rootName)
		},
	}
	return cmd
}

func newRootListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all roots in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := project.AddAllRoots()
			if err != nil {
				return err
			}

			log.Infof("There %s %d %s in the project %s:",
				pluralize("is", "are", len(project.Roots)),
				len(project.Roots),
				pluralize("root", "roots", len(project.Roots)),
				project.ID)
			names := make([]string, len(project.Roots))
			i := 0
			for name := range project.Roots {
				names[i] = name
				i++
			}
			sort.Strings(names)

			for _, name := range names {
				fmt.Println(name)
			}

			return nil
		},
	}
	return cmd
}

// func newRootGraphCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "graph",
// 		Short: "",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			logger.Warn("not yet implemented")
// 			return nil
// 		},
// 	}
// 	return cmd
// }
