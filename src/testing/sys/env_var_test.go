package sys

import (
	"os"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("FOO", "1")
	envVars := GetEnvVars()
	fooVal, found := envVars["FOO"]
	if !found {
		t.Errorf("Could not find environment variable FOO")
	}
	if fooVal != "1" {
		t.Errorf("Env var FOO's value is wrong, expected '1', got %s", fooVal)
	}
}
