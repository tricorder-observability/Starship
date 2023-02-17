package testing

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRunningContainer(t *testing.T) {
	postgresRunner := Runner{
		ImageName: "postgres",
		Options:   []string{"--env=POSTGRES_PASSWORD=passwd"},
		RdyMsg:    "database system is ready to accept connections",
	}

	err := postgresRunner.Launch(10 * time.Second)
	if err != nil {
		t.Fatalf("Could not start postgres, error: %v", err)
	}

	defer func() {
		err := postgresRunner.Stop()
		if err != nil {
			t.Errorf("Failed to stop postgres, error: %v", err)
		}
	}()

	details, err := postgresRunner.Inspect("")
	if err != nil {
		t.Errorf("%s %v", details, err)
	}

	nameSubstr := fmt.Sprintf("\"Name\": \"/%s\",", postgresRunner.ContainerName)
	if !strings.Contains(details, nameSubstr) {
		t.Errorf("Inspect output has no name substr '%s'", nameSubstr)
	}

	_, err = postgresRunner.GetExposedPort(5432)
	if err != nil {
		t.Errorf("Could not obtain exposed port for 5432, error: %s", err)
	}
}
