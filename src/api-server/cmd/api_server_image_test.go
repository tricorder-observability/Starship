package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	testutils "github.com/tricorder/src/testing/bazel"
	docker "github.com/tricorder/src/testing/docker"
)

// Tests that agent image can be executed without missing dynamic library file.
func TestAgentImageHasNoLinkingIssue(t *testing.T) {
	assert := assert.New(t)

	path := "src/api-server/cmd/api-server_image.tar"
	agentPath := testutils.TestBinaryPath(path)

	var cli docker.CLI
	imageName := cli.Load(agentPath)

	runner := docker.Runner{
		ImageName: imageName,
		RdyMsg:    "API server has started",
		AppFlags: []string{
			// Disable management Web UI, so no need to connect to Postgres Grafana
			"--enable_mgmt_ui=false",
			// Disable metadata service, so no need to connect to Postgres
			"--enable_metadata_service=false",
		},
	}
	err := runner.Launch(10 * time.Second)
	defer assert.Nil(runner.Stop())
	assert.Nil(err)
}
