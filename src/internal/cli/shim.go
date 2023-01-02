package cli

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terraboots/internal/shell"
)

func newTerraformCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s ROOT [SCOPE]... [-- TF_FLAG=VALUE]", name),
		Short:   fmt.Sprintf("Runs `terraform %s` in the given root. Pass arguments to terraform after a `--` (for example `terraboots plan ROOT -- -lock=false`)", name),
		Args:    cobra.MinimumNArgs(1),
		GroupID: commandGroupIDTerraform,

		PersistentPreRunE: bootsbootsPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			// log.Infof("args: %+v", args)
			scopes := make([]string, 0, len(args)-1)
			tfargs := make([]string, 0, len(args)-1)
			i := 1
			for i = 1; i < len(args); i++ {
				ok, err := project.IsScopeValue(args[i])
				if err != nil {
					return err
				}
				if ok {
					scopes = append(scopes, args[i])
				} else {
					break
				}
				tfargs = args[i+1:]
			}
			// log.Infof("found scopes: %+v (%d)", scopes, i)
			// log.Infof("remaining args: %+v", args[i:])
			// get a list of locations to run in
			dirs, err := project.BuildRoot(args[0], scopes)
			if err != nil {
				return err
			}

			// TODO: use a worker pool
			for _, dir := range dirs {
				err = runTerraform(name, dir, tfargs, log)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

func runTerraform(command, cwd string, args []string, log *logrus.Entry) error {
	log.Infof("terraform %s in %s with args %+v", command, cwd, args)
	args = append([]string{command}, args...)
	cmd := shell.NewCommand("terraform", args, cwd, log.Logger)

	return cmd.Run()
}
