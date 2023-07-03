// Package cli provides functions relating to running terrascope as a command-
// line interface.
package cli

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/logformatter"
)

var dryRun bool
var quiet bool
var verbose bool
var vertrace bool

var all bool
var yesOne bool
var noNone bool

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

// NewTerrascopeCmd returns a new CLI command representing Terrascope.
func NewTerrascopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terrascope",
		Short: "A build orchestrator for terraform monorepos",
	}

	cmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "don't actually execute the task, just print it out")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase log output")
	cmd.PersistentFlags().BoolVar(&vertrace, "vvv", false, "increase log output even more")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "silences all logs but the errors (and prints those to stderr). Still prints command output to stdout. Overrides verbose and vvv")

	cmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "automatically respond to 'yes/no/all' or 'none/one/all' prompts with 'all'. Overrides both yes and no.")
	cmd.PersistentFlags().BoolVarP(&yesOne, "yes", "y", false, "automatically respond to 'yes/no/all' prompts with 'yes', or 'none/one/all' prompts with 'one'")
	cmd.PersistentFlags().BoolVarP(&noNone, "no", "n", false, "automatically respond to 'yes/no/all' prompts with 'no', or 'none/one/all' prompts with 'none'. Overrides yes.")

	cmd.AddCommand(newVersionCommand())

	cmd.AddCommand(newModuleCommand())
	cmd.AddCommand(newProviderCommand())

	return cmd
}
