package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbose bool

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

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terraboots",
		Short: "A build orchestrator for terraform monorepos",
	}
	return cmd
}
