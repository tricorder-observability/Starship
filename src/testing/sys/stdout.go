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

package sys

import (
	"bytes"
	"io"
	"os"
)

type Stdout struct {
	old *os.File
	r   *os.File
	w   *os.File
}

// Start swap the stdout and allows stdout to be captured.
func (stdout *Stdout) Start() {
	stdout.old = os.Stdout
	stdout.r, stdout.w, _ = os.Pipe()
	os.Stdout = stdout.w
}

// Get restores the original stdout and returns the captured text.
func (stdout *Stdout) Get() string {
	stdout.w.Close()
	os.Stdout = stdout.old

	var buf bytes.Buffer
	_, err := io.Copy(&buf, stdout.r)
	if err != nil {
		return ""
	}
	return buf.String()
}

// CaptureStdout executes the function and capture the stdout during its execution.
// Note that there is no guarantee that the captured stdout actually comes from the function.
func CaptureStdout(fn func()) string {
	stdout := Stdout{}
	stdout.Start()
	fn()
	return stdout.Get()
}
