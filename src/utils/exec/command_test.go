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
