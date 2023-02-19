package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/version"
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.Info()
			fmt.Printf("Terrascope Version: %s\n", info.VersionNumber)
			fmt.Printf("Git Hash: %s\n", info.GitHash)
			fmt.Printf("Build Time: %s\n", info.BuildTime)
			return nil
		},
	}
	return cmd
}
