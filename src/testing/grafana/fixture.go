// Copyright (C) 2023  tricorder-observability
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

package grafana

import (
	"fmt"
	"strings"
	"time"

	docker "github.com/tricorder/src/testing/docker"
)

const defaultPort = 3000

func LaunchContainer() (func() error, string, error) {
	grafanaRunner := &docker.Runner{
		ImageName: "public.ecr.aws/tricorder/grafana:v0.0.9",
		RdyMsg:    "starting",
	}
	err := grafanaRunner.Launch(10 * time.Second)
	if err != nil {
		return nil, "", fmt.Errorf("failed to start postgres server, error: %v", err)
	}

	grafanaGatewayIP, err := grafanaRunner.GetGatewayIP()
	if err != nil {
		return nil, "", err
	}

	pgPort, err := grafanaRunner.GetExposedPort(defaultPort)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get grafana container's exposed port for %d, error: %v", defaultPort, err)
	}
	grafanaURL := fmt.Sprintf("http://%s:%d", strings.TrimSpace(grafanaGatewayIP), pgPort)

	cleanerFn := func() error {
		if err := grafanaRunner.Stop(); err != nil {
			return err
		}
		return nil
	}
	return cleanerFn, grafanaURL, nil
}
