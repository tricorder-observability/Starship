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

package file

import (
	"os"
	"path"
	"testing"

	// Use primitive bazel instead of src/testing/bazel to avoid circular dependency between
	// src/utils/fail and src/testing/bazel
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testdataGoodWASM = "src/utils/file/testdata/good.wasm"
	testdataBadWASM1 = "src/utils/file/testdata/bad_fmt.wasm"
	testdataBadWASM2 = "src/utils/file/testdata/bad_magic_num.wasm"
	testdataBadWASM3 = "src/utils/file/testdata/bad_suffix.wa"
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

// Tests IsWasmELF() can check if a file is a valid wasm file.
func TestIsWasmELF(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	goodWASM, err := runfiles.Rlocation("tricorder/" + testdataGoodWASM)
	require.NoError(err)
	badWASM1, err := runfiles.Rlocation("tricorder/" + testdataBadWASM1)
	require.NoError(err)
	badWASM2, err := runfiles.Rlocation("tricorder/" + testdataBadWASM2)
	require.NoError(err)
	badWASM3, err := runfiles.Rlocation("tricorder/" + testdataBadWASM3)
	require.NoError(err)

	assert.True(IsWasmELF(goodWASM))
	assert.False(IsWasmELF(badWASM1))
	assert.False(IsWasmELF(badWASM2))
	assert.False(IsWasmELF(badWASM3))
}

// Tests GetFileType() returns correct file type consts.
func TestGetFileType(t *testing.T) {
	assert := assert.New(t)

	for _, c := range []struct {
		filePath  string
		expedType string
	}{
		{"test.wasm", WASM},
		{"test.c", C},
		{"test.bcc", BCC},
		{"test.unknown", UNKNOWN},
	} {
		assert.Equal(GetFileType(c.filePath), c.expedType)
	}
}
