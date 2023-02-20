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

package pg

import (
	"fmt"
	"strings"
	"time"

	testutils "github.com/tricorder/src/testing/bazel"
	docker "github.com/tricorder/src/testing/docker"
	pgutils "github.com/tricorder/src/utils/pg"
)

const (
	defaultPort       = 5432
	path              = "src/testing/pg/postgres.tar"
	postgresImageName = "bazel/src/testing/pg:postgres"
)

// init loads the pulled postgres container image when importing this package.
// So the postgres image only needs to be pulled once during one bazel command run.
func init() {
	var cli docker.CLI
	cli.Load(testutils.TestFilePath(path))
}

// LaunchContainer returns a cleaner function, the client connected to the started container, and error if failed.
// You can then destroy the fixtures by deferring statement:
// cleaner, pgClient, err := createPGTestFixutre()
// require.Nil(err)
//
//	defer func() {
//	  assert.Nil(cleaner())
//	}()
func LaunchContainer() (func() error, *pgutils.Client, error) {
	pgRunner := &docker.Runner{
		ImageName: postgresImageName,
		RdyMsg:    "database system is ready to accept connections",
		Options:   []string{"--env=POSTGRES_PASSWORD=passwd"},
	}
	err := pgRunner.Launch(10 * time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to start postgres server, error: %v", err)
	}

	pgGatewayIP, err := pgRunner.GetGatewayIP()
	if err != nil {
		return nil, nil, err
	}

	pgPort, err := pgRunner.GetExposedPort(defaultPort)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get postgres container's exposed port for %d, error: %v",
			defaultPort,
			err,
		)
	}
	pgURL := fmt.Sprintf("postgresql://postgres:passwd@%s:%d", strings.TrimSpace(pgGatewayIP), pgPort)

	pgClient := pgutils.NewClient(pgURL)
	if err := pgClient.Connect(); err != nil {
		return nil, nil, fmt.Errorf("Unable to create client to postgres at %s, error: %v", pgURL, err)
	}

	cleanerFn := func() error {
		pgClient.Close()
		if err := pgRunner.Stop(); err != nil {
			return err
		}
		return nil
	}
	return cleanerFn, pgClient, nil
}
