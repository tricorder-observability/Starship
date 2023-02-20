// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
