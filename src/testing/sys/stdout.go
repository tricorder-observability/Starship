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
