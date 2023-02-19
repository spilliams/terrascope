package cli

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/hcl"
)

func newGraphModuleCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "graph [DIRECTORY]",
		Short:   "graphs the root module at the given directory (`.` by default)",
		GroupID: commandGroupIDTunnelvision,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir, err := os.Getwd()
			if err != nil {
				return err
			}
			if len(args) > 0 {
				rootDir = args[0]
				rootDir, err = filepath.Abs(rootDir)
				if err != nil {
					return err
				}
			}
			// outFilename := "output.dot"

			logrus.Infof("reading configuration at %s", rootDir)
			// logrus.Infof("outputting graph in file %s", outFilename)

			parser := hcl.NewModule(logrus.StandardLogger())
			parser.ParseModuleDirectory(rootDir)

			logrus.Debugf("%#v", parser.Parser())
			logrus.Debugf("%#v", parser.Module())

			// TODO: build a graph from the module

			return nil
		},
	}
}
