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

	path := "src/agent/cmd/agent_image.tar"
	agentPath := testutils.TestBinaryPath(path)

	var cli docker.CLI
	imageName := cli.Load(agentPath)

	agentRunner := docker.Runner{
		ImageName: imageName,
		RdyMsg:    "connecting to API server",
	}
	assert.Nil(agentRunner.Launch(5 * time.Second))
}
