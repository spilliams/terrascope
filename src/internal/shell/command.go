// Package shell wraps exec.Cmd with some conveniences for capturing output
package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Command wraps an `exec.Cmd`, with extra logging goodies
type Command struct {
	shellCmd *exec.Cmd

	*logrus.Entry
}

// NewCommand builds a new `Command`, given basic structure and a logger
func NewCommand(command string, args []string, workingDir string, logger *logrus.Logger) *Command {
	shellCmd := exec.Command(command, args...)
	shellCmd.Dir = workingDir
	shellCmd.Stdin = os.Stdin

	return &Command{
		shellCmd: shellCmd,
		Entry:    logger.WithField("prefix", fmt.Sprintf("tf %s %s", args[0], workingDir)),
	}
}

func (c *Command) String() string {
	return fmt.Sprintf("cd %s && %s", c.shellCmd.Dir, c.shellCmd.String())
}

// Run runs the receiver, with stdout piped to the logger's debug Level, and
// stderr piped to the logger's Warn level
func (c *Command) Run() error {
	c.Debugf(c.shellCmd.String())

	stdout, err := c.shellCmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := c.shellCmd.StderrPipe()
	if err != nil {
		return err
	}

	err = c.shellCmd.Start()
	if err != nil {
		return err
	}

	err = readStdoutAndStderr(c.Entry, stdout, stderr)
	if err != nil {
		return err
	}

	return c.shellCmd.Wait()
}
