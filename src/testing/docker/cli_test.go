package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testutils "github.com/tricorder/src/testing/bazel"
)

// Tests that Load() can load a .tar file.
func TestLoad(t *testing.T) {
	assert := assert.New(t)

	var cli CLI
	tarFilePath := testutils.TestFilePath("src/testing/docker/testdata/test_image.tar")
	assert.Equal("bazel/src/testing/docker/testdata:test_image", cli.Load(tarFilePath))
}
