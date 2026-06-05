package cli

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/shell"
)

func newSpecificTerraformCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s ROOT [SCOPE]... [-- TF_FLAG=VALUE]", name),
		Short:   fmt.Sprintf("Runs `terraform %s` in the given root", name),
		Long:    fmt.Sprintf("Runs `terraform %s` in the given root. Pass arguments to terraform after a `--` (for example `terrascope %s ROOT -- -lock=false`)", name, name),
		Args:    cobra.MinimumNArgs(1),
		GroupID: commandGroupIDTerraformShim,

		PersistentPreRunE: parseProject,
		RunE: func(cmd *cobra.Command, args []string) error {
			tf, err := buildTerraformCommand(args, name)
			if err != nil {
				return err
			}

			tf.run()
			return nil
		},
	}

	return cmd
}

func newGenericTerraformCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("tf ROOT [SCOPE]... -- COMMAND [TF_FLAG=VALUE]..."),
		Aliases: []string{"terraform"},
		Short:   "Runs a given terraform command in the given root",
		Long:    fmt.Sprintf("Runs a given terraform command in the given root. Pass arguments to terraform after a `--` (for example `terrascope tf ROOT -- state list`)"),
		Args:    cobra.MinimumNArgs(1),
		GroupID: commandGroupIDTerraformShim,

		PersistentPreRunE: parseProject,
		RunE: func(cmd *cobra.Command, args []string) error {
			tf, err := buildTerraformCommand(args, "")
			if err != nil {
				return err
			}

			tf.run()
			return nil
		},
	}

	return cmd
}

type tfCmd struct {
	dirs []string
	args []string
}

func buildTerraformCommand(args []string, cmd string) (*tfCmd, error) {
	err := project.AddAllRoots()
	if err != nil {
		return nil, err
	}

	// log.Infof("args: %+v", args)
	scopes := make([]string, 0, len(args)-1)
	tfargs := make([]string, 0, len(args)-1)
	if len(cmd) > 0 {
		tfargs = append(tfargs, cmd)
	}
	i := 1
	for i = 1; i < len(args); i++ {
		ok, err := project.IsScopeValue(args[i])
		if err != nil {
			return nil, err
		}
		if ok {
			scopes = append(scopes, args[i])
		} else {
			break
		}
	}
	tfargs = append(tfargs, args[i:]...)
	log.Infof("found scopes: %+v", scopes)
	log.Infof("remaining args: %+v", args[i:])
	// get a list of locations to run in
	dirs, err := project.BuildRoot(args[0], scopes, dryRun, chainDependenciesOption())
	if err != nil {
		return nil, err
	}

	return &tfCmd{
		dirs: dirs,
		args: tfargs,
	}, nil
}

func (tf *tfCmd) run() {
	var wg sync.WaitGroup
	wg.Add(len(tf.dirs))
	for _, dir := range tf.dirs {
		go runTerraform(&wg, dir, tf.args, log)
	}
	wg.Wait()
}

func runTerraform(wg *sync.WaitGroup, cwd string, args []string, log *logrus.Entry) {
	cmd := shell.NewCommand("terraform", args, cwd, log.Logger)
	if dryRun {
		fmt.Printf("%s\n", cmd.String())
		err := cmd.Run()
		if err != nil {
			log.Error(err.Error())
		}
	}
	wg.Done()
}
