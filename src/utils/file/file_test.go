package file

import (
	"os"
	"path"
	"testing"

	// Use primitive bazel instead of src/testing/bazel to avoid circular dependency between
	// src/utils/fail and src/testing/bazel
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	if Exists("non-exitent") {
		t.Errorf("File should not exist")
	}

	path := bazel.TestTmpDir()
	if !Exists(path) {
		t.Errorf("File '%s' should exist", path)
	}
}

// Tests that Create() can create a file.
func TestCreateAppend(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "a/b/c/d/e")
	assert.Nil(Create(p))

	assert.Nil(Append(p, "test"))
	content, err := os.ReadFile(p)
	assert.Nil(err)
	assert.Equal("test", string(content))

	assert.Nil(Append(p, "foo"))
	content, err = os.ReadFile(p)
	assert.Nil(err)
	assert.Equal("testfoo", string(content))
}

// Tests that TestCreateDir().
func TestCreateDir(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := tmpDir + "a/b/c/d/e"
	assert.Nil(CreateDir(p))

	assert.Equal(Exists(p), true)
	s, err := os.Stat(p)
	assert.Nil(err)
	assert.Equal(s.IsDir(), true)
}

// Tests ReadLines() can read file with multiple lines into string slice.
func TestReadlines(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")
	assert.Nil(Write(p, "123\n456\n789"))

	results, err := ReadLines(p)
	assert.Nil(err)
	assert.Equal([]string{"123", "456", "789"}, results)
}

// Tests ReadLines() can read file with multiple lines into string slice.
func TestReadBin(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")
	assert.Nil(Write(p, "123\n456\n789"))

	results, err := ReadBin(p)
	assert.Nil(err)
	assert.Equal("123\n456\n789", string(results))
}

// Tests Copy() can copy srcPath file to dstPath file.
func TestCopy(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")
	assert.Nil(Write(p, "123\n456\n789"))

	newTmpDir := bazel.TestTmpDir()
	np := path.Join(newTmpDir, "readlines_test_1")
	assert.Nil(Copy(p, np))
	results, err := ReadBin(np)
	assert.Nil(err)
	assert.Equal("123\n456\n789", string(results))
}

// Tests Reader() can return reader and closer object.
func TestReader(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")
	assert.Nil(Write(p, "123\n456\n789"))

	reader, closer, err := Reader(p)
	assert.Nil(err)
	a := make([]byte, 12)
	results, err := reader.Read(a)
	assert.Nil(err)
	assert.Equal(results, 11)
	assert.Equal(string(a), "123\n456\n789\x00")
	assert.Nil(closer.Close())
	_, err = reader.Read(a)
	assert.NotNil(err)
}

// Tests Writer() can return writer and closer object.
func TestWriter(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")

	writer, closer, err := Writer(p)
	assert.Nil(err)
	results, err := writer.Write([]byte("123\n456\n789"))
	assert.Nil(err)
	assert.Nil(closer.Close())
	assert.Equal(results, 11)

	a := make([]byte, 12)
	reader, closer, err := Reader(p)
	assert.Nil(err)
	results, err = reader.Read(a)
	assert.Nil(err)
	assert.Equal(results, 11)
	assert.Equal(string(a), "123\n456\n789\x00")
	assert.Nil(closer.Close())
	_, err = writer.Write([]byte("123\n456\n789"))
	assert.NotNil(err)
}

// Tests ReadSymlink() can read symlink file.
func TestReadSymlink(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")

	writer, closer, err := Writer(p)
	assert.Nil(err)
	results, err := writer.Write([]byte("123\n456\n789"))
	assert.Nil(err)
	assert.Nil(closer.Close())
	assert.Equal(results, 11)

	np := tmpDir + "/a/b/c/e"
	assert.Nil(CreateSymLink(p, np))
	link, err := ReadSymLink(np)
	assert.Nil(err)
	assert.Equal(link, p)
}

// Tests CreateSymLink() can copy srcPath file to dstPath file.
func TestCreateSymLink(t *testing.T) {
	assert := assert.New(t)

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")

	writer, closer, err := Writer(p)
	assert.Nil(err)
	results, err := writer.Write([]byte("123\n456\n789"))
	assert.Nil(err)
	assert.Nil(closer.Close())
	assert.Equal(results, 11)

	np := tmpDir + "/a/b/c/g"
	assert.Nil(CreateSymLink(p, np))

	a := make([]byte, 12)
	reader, closer, err := Reader(np)
	assert.Nil(err)
	results, err = reader.Read(a)
	assert.Nil(err)
	assert.Equal(results, 11)
	assert.Equal(string(a), "123\n456\n789\x00")
	assert.Nil(closer.Close())
}

// Tests Contains() can check if a file contains a string.
func TestContains(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(false, Contains("non-exist-file", "bar"))

	tmpDir := bazel.TestTmpDir()
	p := path.Join(tmpDir, "readlines_test")
	assert.Nil(Write(p, "123\n456\n789"))

	assert.Equal(false, Contains(p, "bar"))
	assert.Equal(true, Contains(p, "123"))
}
