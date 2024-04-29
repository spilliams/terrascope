package cli

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/pkg/terrascope"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "root",
		Aliases: []string{"r"},
		Short:   "Commands relating to root modules",
		GroupID: commandGroupIDTerrascope,

		PersistentPreRunE: parseProject,
	}

	cmd.AddCommand(newRootBuildCommand())
	cmd.AddCommand(newRootCleanCommand())
	cmd.AddCommand(newRootGenerateCommand())
	cmd.AddCommand(newRootGraphResourcesCommand())
	cmd.AddCommand(newRootListCommand())
	cmd.AddCommand(newRootShowCommand())

	return cmd
}

func newRootBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build ROOT [SCOPE]...",
		Aliases: []string{"b"},
		Short:   "Builds the given root and prints the location of the built configurations to stdout",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := project.AddAllRoots()
			if err != nil {
				return err
			}

			scopes := make([]string, len(args)-1)
			for i := 1; i < len(args); i++ {
				scopes[i-1] = args[i]
			}
			dirs, err := project.BuildRoot(args[0], scopes, dryRun, chainDependenciesOption())

			for _, dir := range dirs {
				fmt.Println(dir)
			}
			return err
		},
	}

	return cmd
}

func newRootCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean ROOT [SCOPE]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Cleans a root of all generated files",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := project.AddAllRoots()
			if err != nil {
				return err
			}

			scopes := make([]string, len(args)-1)
			for i := 1; i < len(args); i++ {
				scopes[i-1] = args[i]
			}
			dirs, err := project.CleanRoot(args[0], scopes, dryRun)
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

func newRootGraphResourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph-resources ROOT",
		Short: "(EXPERIMENTAL) Prints out a DOT-format graph of the resource and data blocks in the given root.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := project.AddAllRoots()
			if err != nil {
				return err
			}
			root := project.GetRoot(args[0])
			return printModuleGraph(path.Dir(root.Filename))
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

func newRootShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show ROOT",
		Short: "Prints information about a root",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := project.AddAllRoots()
			if err != nil {
				return err
			}

			root := project.GetRoot(args[0])
			if root == nil {
				log.Warnf("No root named %s was found", args[0])
			}

			file, err := ioutil.ReadFile(root.Filename)
			if err != nil {
				return err
			}
			fmt.Println(string(file))
			return nil
		},
	}

	return cmd
}

func chainDependenciesOption() terrascope.RootDependencyChain {
	if all {
		return terrascope.RootDependencyChainAll
	}
	if noNone {
		return terrascope.RootDependencyChainNone
	}
	if yesOne {
		return terrascope.RootDependencyChainOne
	}
	return terrascope.RootDependencyChainUnknown
}
