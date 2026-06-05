package cli

import (
	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/pkg/server"
)

func newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.Run()
		},
	}
}
