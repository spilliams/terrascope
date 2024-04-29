// Package cli provides functions relating to running terrascope as a command-
// line interface.
package cli

import (
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/logformatter"
	"github.com/spilliams/terrascope/pkg/terrascope"
)

var configFile string
var dryRun bool
var quiet bool
var verbose bool
var vertrace bool

var log *logrus.Entry

func init() {
	cobra.OnInitialize(initLogger)
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
	if vertrace {
		logger.SetLevel(logrus.TraceLevel)
	}
	if quiet {
		logger.SetLevel(logrus.ErrorLevel)
		logger.SetOutput(os.Stderr)
	}
	log = logger.WithField("prefix", "main")
}

const commandGroupIDTerraformShim = "terraform-shim"
const commandGroupIDTerrascope = "terrascope"
const commandGroupIDTerraformTools = "terraform-tools"

var project *terrascope.Project

// NewTerrascopeCmd returns a new CLI command representing Terrascope.
func NewTerrascopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terrascope",
		Short: "A build orchestrator for terraform monorepos",
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase log output")
	cmd.PersistentFlags().BoolVar(&vertrace, "vvv", false, "increase log output even more")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "silences all logs but the errors (and prints those to stderr). Still prints command output to stdout. Overrides verbose and vvv")

	cmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "automatically respond to 'yes/no/all' or 'none/one/all' prompts with 'all'. Overrides both yes and no.")
	cmd.PersistentFlags().BoolVarP(&yesOne, "yes", "y", false, "automatically respond to 'yes/no/all' prompts with 'yes', or 'none/one/all' prompts with 'one'")
	cmd.PersistentFlags().BoolVarP(&noNone, "no", "n", false, "automatically respond to 'yes/no/all' prompts with 'no', or 'none/one/all' prompts with 'none'. Overrides yes.")

	cmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "terrascope.hcl", "the filename of the project configuration")

	cmd.AddCommand(newVersionCommand())

	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerrascope, Title: "Working with your terrascope project"})
	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerraformShim, Title: "Terraform Commands"})
	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerraformTools, Title: "Inspecting your infrastructure"})

	cmd.AddCommand(newSpecificTerraformCommand("init"))
	cmd.AddCommand(newSpecificTerraformCommand("plan"))
	cmd.AddCommand(newSpecificTerraformCommand("apply"))
	cmd.AddCommand(newSpecificTerraformCommand("destroy"))
	cmd.AddCommand(newSpecificTerraformCommand("output"))
	cmd.AddCommand(newSpecificTerraformCommand("console"))
	cmd.AddCommand(newGenericTerraformCommand())

	cmd.AddCommand(newProjectCommand())
	cmd.AddCommand(newScopeCommand())
	cmd.AddCommand(newRootCommand())

	cmd.AddCommand(newModuleCommand())
	cmd.AddCommand(newProviderCommand())

	return cmd
}

func parseProject(cmd *cobra.Command, args []string) error {
	log.Debugf("Using project configuration file: %s", configFile)
	var err error
	project, err = terrascope.ParseProject(configFile, log.Logger)
	if err != nil {
		return err
	}

	rootsDir := path.Join(path.Dir(configFile), project.RootsDir)
	rootsDir, err = filepath.Abs(rootsDir)
	if err != nil {
		return err
	}
	project.RootsDir = rootsDir
	// log.Debugf("Project roots directory: %s", project.RootsDir)
	// log.Debugf("Project scope data files: %s", project.ScopeDataFiles)

	return nil
}
