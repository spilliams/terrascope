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

var quiet bool
var verbose bool
var vertrace bool
var configFile string
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

const commandGroupIDTerraform = "terraform"
const commandGroupIDTerrascope = "terrascope"
const commandGroupIDTunnelvision = "tunnelvision"

var project *terrascope.Project

func NewTerrascopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terrascope",
		Short: "A build orchestrator for terraform monorepos",
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase log output")
	cmd.PersistentFlags().BoolVar(&vertrace, "vvv", false, "increase log output even more")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "silences all logs but the errors (and prints those to stderr). Still prints command output to stdout. Overrides verbose and vvv")

	cmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "terrascope.hcl", "the filename of the project configuration")

	// TODO: version command
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerrascope, Title: "Working with your terrascope project"})
	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTerraform, Title: "Terraform Commands"})
	cmd.AddGroup(&cobra.Group{ID: commandGroupIDTunnelvision, Title: "Inspecting your infrastructure"})

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
