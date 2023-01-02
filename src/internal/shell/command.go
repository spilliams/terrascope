package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Command struct {
	command    string
	args       []string
	workingDir string
	*logrus.Entry
}

func NewCommand(command string, args []string, workingDir string, logger *logrus.Logger) *Command {
	return &Command{
		command:    command,
		args:       args,
		workingDir: workingDir,
		Entry:      logger.WithField("prefix", fmt.Sprintf("tf %s %s", args[0], workingDir)),
	}
}

func (c *Command) Run() error {
	shellCmd := exec.Command(c.command, c.args...)
	shellCmd.Dir = c.workingDir
	shellCmd.Stdin = os.Stdin
	c.Debugf(shellCmd.String())

	stdout, err := shellCmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := shellCmd.StderrPipe()
	if err != nil {
		return err
	}

	err = shellCmd.Start()
	if err != nil {
		return err
	}

	err = readStdoutAndStderr(c.Entry, stdout, stderr)
	if err != nil {
		return err
	}

	return shellCmd.Wait()
}
