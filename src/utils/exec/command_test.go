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

import "testing"

func TestCommand(t *testing.T) {
	t.Log("Testing Command APIs")

	argv := []string{"echo", "-n", "hello world"}
	cmd := NewCommand(argv)
	err := cmd.Start()
	if err != nil {
		t.Errorf("Could not start command %v, error: %v", argv, err)
	}
	err = cmd.Wait()
	if err != nil {
		t.Errorf("Could not wait command %v, error: %v", argv, err)
	}
	if cmd.Stderr() != "" {
		t.Errorf("Stderr should be empty, got '%s'", cmd.Stderr())
	}
	expStdout := "hello world"
	if cmd.Stdout() != expStdout {
		t.Errorf("Stdout should be '%s', got '%s'", expStdout, cmd.Stdout())
	}
}
