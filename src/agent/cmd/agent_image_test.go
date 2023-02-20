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
