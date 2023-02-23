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

package testing

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunningContainer(t *testing.T) {
	assert := assert.New(t)

	postgresRunner := Runner{
		ImageName: "postgres",
		EnvVars:   map[string]string{"a": "b"},
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

	out := postgresRunner.Exec([]string{"bash", "-c", "echo -n ${POSTGRES_PASSWORD}"})
	assert.Equal("passwd", out)
}
