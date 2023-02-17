package exec

import (
	"bytes"
	"os/exec"
)

type Command struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	cmd    *exec.Cmd
}

func NewCommand(argv []string) *Command {
	cmd := new(Command)

	cmd.stdout = bytes.Buffer{}
	cmd.stderr = bytes.Buffer{}

	cmd.cmd = exec.Command(argv[0], argv[1:]...)
	cmd.cmd.Stdout = &cmd.stdout
	cmd.cmd.Stderr = &cmd.stderr

	return cmd
}

func (cmd *Command) Start() error {
	return cmd.cmd.Start()
}

func (cmd *Command) Wait() error {
	return cmd.cmd.Wait()
}

func (cmd *Command) Stdout() string {
	return cmd.stdout.String()
}

func (cmd *Command) Stderr() string {
	return cmd.stderr.String()
}
