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
