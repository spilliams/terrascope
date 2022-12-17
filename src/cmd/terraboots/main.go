package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terraboots/internal/terraboots"
)

var verbose bool
var configFile string

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
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
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

	cmd.AddCommand(newTerraformCommand("init"))
	cmd.AddCommand(newTerraformCommand("plan"))
	cmd.AddCommand(newTerraformCommand("apply"))
	cmd.AddCommand(newTerraformCommand("destroy"))
	cmd.AddCommand(newTerraformCommand("output"))
	cmd.AddCommand(newTerraformCommand("console"))

	// TODO: manual tf command
	// cmd.AddCommand(newManualTfCommand())

	cmd.AddCommand(newScopeCommand())
	cmd.AddCommand(newRootCommand())

	// TODO: command to generate a new project
	// cmd.AddCommand(newProjectCommand())

	return cmd
}

func newTerraformCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s ROOT", name),
		Short:   fmt.Sprintf("Runs `terraform %s` in the given root", name),
		GroupID: commandGroupIDTerraform,
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

func newScopeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scope",
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
		Use:   "list",
		Short: "Lists all scopes in this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

func newScopeGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new scope in this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "root",
		Short:   "Commands relating to root modules",
		GroupID: commandGroupIDTerraboots,

		PersistentPreRunE: bootsbootsPreRunE,
	}

	cmd.AddCommand(newRootBuildCommand())
	cmd.AddCommand(newRootGenerateCommand())
	cmd.AddCommand(newRootGraphCommand())
	cmd.AddCommand(newRootListCommand())

	return cmd
}

func newRootBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build ROOT",
		Short: "Builds the given root and prints the location of the built root to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

func newRootGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

func newRootGraphCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

func newRootListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Warn("not yet implemented")
			return nil
		},
	}

	return cmd
}

var project *terraboots.Project

func bootsbootsPreRunE(cmd *cobra.Command, args []string) error {
	logrus.Debug("hi i'm bootsboots")

	var err error
	project, err = terraboots.ParseProject(configFile)
	if err != nil {
		return err
	}

	return nil
}
