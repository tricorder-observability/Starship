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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/file"
)

const kprobeEventsContent = `p:kprobes/p___x64_sys_read_bcc_212599 __x64_sys_read
e:kprobes/p___x64_sys_read_bcc_212599 __x64_sys_read
p:kprobes/p___x64_sys_read_bcc_212599 __x64_sys_read extra_errorneous_filed`

// Tests that findProbes() find the probes with the marker.
func TestFindProbes(t *testing.T) {
	assert := assert.New(t)

	kprobeFile := testutils.CreateTmpFileWithContent(kprobeEventsContent)

	probes, err := findProbes(kprobeFile, "_bcc_")
	assert.Nil(err)
	assert.Equal([]string{"p:kprobes/p___x64_sys_read_bcc_212599"}, probes)
}

// Tests that cleanProbes() writes the delete probe lines to the specified file.
func TestCleanProbes(t *testing.T) {
	assert := assert.New(t)

	kprobeFile := testutils.CreateTmpFile()

	err := cleanProbes(kprobeFile, []string{":a", ":b", ":c"})
	assert.Nil(err)
	content, err := file.Read(kprobeFile)
	assert.Nil(err)
	assert.Equal("-:a\n-:b\n-:c", content)
}
