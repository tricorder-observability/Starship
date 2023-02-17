package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/utils/file"
)

func TestCreateTmpDir(t *testing.T) {
	assert := assert.New(t)

	dir1 := CreateTmpDir()
	dir2 := CreateTmpDir()
	assert.NotEqual(dir1, dir2)
	assert.Contains(dir1, "tricorder-")
	assert.Contains(dir2, "tricorder-")
}

func TestCreateTmpFile(t *testing.T) {
	assert := assert.New(t)

	file1 := CreateTmpFile()
	assert.True(file.Exists(file1))

	file2 := CreateTmpFile()
	assert.True(file.Exists(file2))

	assert.NotEqual(file1, file2)
}

// Tests that ReadTestFile() and ReadTestBinFile() work as expected.
func TestReadTestFile(t *testing.T) {
	assert := assert.New(t)
	content, err := ReadTestFile("src/testing/bazel/test")
	assert.Nil(err)
	assert.Equal("hello test\n", content)
	binContent, err := ReadTestBinFile("src/testing/bazel/test")
	assert.Nil(err)
	assert.Equal([]byte("hello test\n"), binContent)
}
