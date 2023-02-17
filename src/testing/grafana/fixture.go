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
