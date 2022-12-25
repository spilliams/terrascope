package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terraboots/internal/terraboots"
	"github.com/spilliams/terraboots/pkg/logformatter"
)

var verbose bool
var configFile string
var log *logrus.Entry

func init() {
	cobra.OnInitialize(initLogger)
}

func main() {
	err := newRootCmd().Execute()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func initLogger() {
	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logformatter.PrefixedTextFormatter{
		UseColor: true,
	})
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}
	log = logger.WithField("prefix", "main")
}

const commandGroupIDTerraform = "terraform"
const commandGroupIDTerraboots = "terraboots"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terraboots",
		Short: "A build orchestrator for terraform monorepos",
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase log output")
	cmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "terraboots.hcl", "the filename of the project configuration")

	// version
	// help
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerraboots, Title: "Working with your terraboots project"})
	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerraform, Title: "Terraform Commands"})

	// cmd.AddCommand(newTerraformCommand("init"))
	// cmd.AddCommand(newTerraformCommand("plan"))
	// cmd.AddCommand(newTerraformCommand("apply"))
	// cmd.AddCommand(newTerraformCommand("destroy"))
	// cmd.AddCommand(newTerraformCommand("output"))
	// cmd.AddCommand(newTerraformCommand("console"))

	cmd.AddCommand(newScopeCommand())
	cmd.AddCommand(newRootCommand())

	return cmd
}

// func newTerraformCommand(name string) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     fmt.Sprintf("%s ROOT", name),
// 		Short:   fmt.Sprintf("Runs `terraform %s` in the given root", name),
// 		GroupID: commandGroupIDTerraform,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			logger.Warn("not yet implemented")
// 			return nil
// 		},
// 	}

// 	return cmd
// }

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

			log.Infof("There are %d scopes in the project %s", len(project.ScopeTypes), project.ID)
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

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "root",
		Aliases: []string{"r"},
		Short:   "Commands relating to root modules",
		GroupID: commandGroupIDTerraboots,

		PersistentPreRunE: bootsbootsPreRunE,
	}

	cmd.AddCommand(newRootBuildCommand())
	// cmd.AddCommand(newRootGenerateCommand())
	// cmd.AddCommand(newRootGraphCommand())
	// cmd.AddCommand(newRootListCommand())

	return cmd
}

func newRootBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build ROOT",
		Aliases: []string{"b"},
		Short:   "Builds the given root and prints the location of the built root to stdout",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := project.BuildRoot(args[0])
			log.Debugf("%+v\n", root)
			return err
		},
	}

	return cmd
}

// func newRootGenerateCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "generate",
// 		Short: "",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			logger.Warn("not yet implemented")
// 			return nil
// 		},
// 	}

// 	return cmd
// }

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

// func newRootListCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "list",

// 		Short: "",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			logger.Warn("not yet implemented")
// 			return nil
// 		},
// 	}

// 	return cmd
// }

var project *terraboots.Project

func bootsbootsPreRunE(cmd *cobra.Command, args []string) error {
	log.Debugf("Using project configuration file: %s", configFile)
	var err error
	project, err = terraboots.ParseProject(configFile, log.Logger)
	if err != nil {
		return err
	}

	rootsDir := path.Join(path.Dir(configFile), project.RootsDir)
	rootsDir, err = filepath.Abs(rootsDir)
	if err != nil {
		return err
	}
	project.RootsDir = rootsDir
	logrus.Debugf("Project roots directory: %s", project.RootsDir)

	logrus.Debugf("Project scope data files: %s", project.ScopeDataFiles)

	return nil
}
