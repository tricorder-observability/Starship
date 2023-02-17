package tar

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	testuitls "github.com/tricorder/src/testing/bazel"

	"github.com/tricorder/src/utils/file"
)

func TestGZExtract(t *testing.T) {
	assert := assert.New(t)

	tmpDir := testuitls.CreateTmpDir()
	assert.Nil(GZExtract("testdata/test.tar.gz", tmpDir))
	helloPath := path.Join(tmpDir, "hello.txt")
	assert.Equal(file.Exists(helloPath), true)

	tmpDir = testuitls.CreateTmpDir()
	assert.NotNil(GZExtract("testdata/wrong_file_format.tar.gz", tmpDir))
}
