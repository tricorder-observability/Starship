package linux_headers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/file"
)

func TestDistance(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(0, distance(Version{4, 8, 0}, Version{4, 8, 0}))
	assert.Equal(1, distance(Version{4, 8, 1}, Version{4, 8, 0}))
	assert.Equal(101, distance(Version{4, 9, 1}, Version{4, 8, 0}))
	assert.Equal(10101, distance(Version{5, 9, 1}, Version{4, 8, 0}))
}

// Tests that unameVersion() returns the correct release version string.
func TestUnameVersion(t *testing.T) {
	assert := assert.New(t)
	_, err := unameVersion()
	assert.Nil(err)
}

// Tests that unameStr() returns the correct release version string.
func TestUnameStr(t *testing.T) {
	assert := assert.New(t)
	_, err := unameStr()
	assert.Nil(err)
}

// Tests that procVersion() no error.
func TestProcVersion(t *testing.T) {
	assert := assert.New(t)
	_, err := procVersion()
	assert.Nil(err)
}

// Tests that parseVersion() returns the correct release version.
func TestParseVersion(t *testing.T) {
	assert := assert.New(t)

	for _, c := range []struct {
		str string
		v   Version
	}{
		{
			"5.4.228",
			Version{5, 4, 228},
		},
		{
			"1.2.3adfadf",
			Version{1, 2, 3},
		},
		{
			"5.15.0-1028.32-aws",
			Version{5, 15, 0},
		},
	} {
		v, err := parseVersion(c.str)
		assert.Nil(err)
		assert.Equal(c.v, v)
	}

	for _, str := range []string{
		"aaa-1.2.3",
		"1.2.",
	} {
		_, err := parseVersion(str)
		assert.NotNil(err)
	}
}

// Tests that the Veresion.code() returns the right result.
func TestVersionCode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(uint32(266753), Version{4, 18, 1}.code())
}

// Tests that WriteVersion() writes the correct content to the file.
func TestWriteVersion(t *testing.T) {
	assert := assert.New(t)

	p := testutils.CreateTmpFileWithContent(
		"#define LINUX_VERSION_CODE 000000\n" +
			"#define KERNEL_VERSION(a,b,c) (((a) << 16) + ((b) << 8) + (c))")

	assert.Nil(WriteVersion(p, Version{4, 18, 1}))

	c, err := file.Read(p)
	assert.Nil(err)
	assert.Equal(
		"#define LINUX_VERSION_CODE 266753\n"+
			"#define KERNEL_VERSION(a,b,c) (((a) << 16) + ((b) << 8) + (c))", c)
}

func TestSemVerString(t *testing.T) {
	assert := assert.New(t)
	for _, c := range []struct {
		v      Version
		semVer string
	}{
		{Version{1, 0, 0}, "1.0.0"},
		{Version{1, 1, 0}, "1.1.0"},
	} {
		assert.Equal(c.semVer, c.v.semVerStr())
	}
}

// Tests that getKernelVersionFromArchiveFilePath() return the correct version.
func TestGetKernelVersionFromArchiveFilePath(t *testing.T) {
	assert := assert.New(t)

	for _, c := range []struct {
		str string
		v   Version
	}{
		{
			"linux-headers-1.2.3.tar.gz",
			Version{1, 2, 3},
		},
		{
			"linux-headers-1.2.3adfadf.tar.gz",
			Version{1, 2, 3},
		},
		{
			"linux-headers-5.15.0-1028.32-aws.tar.gz",
			Version{5, 15, 0},
		},
	} {
		v, err := getKernelVersionFromArchiveFilePath(c.str)
		assert.Nil(err)
		assert.Equal(c.v, v)
	}

	for _, str := range []string{
		"aaa-1.2.3",
		"1.2.",
	} {
		_, err := getKernelVersionFromArchiveFilePath(str)
		assert.NotNil(err)
	}

	for _, str := range []string{
		"linux-headers-aaa-1.2.3",
		"linux-headers-1.2.",
	} {
		_, err := getKernelVersionFromArchiveFilePath(str)
		assert.NotNil(err)
	}

	for _, str := range []string{
		"aaa-1.2.3",
		"1.2.tar.gz",
	} {
		_, err := getKernelVersionFromArchiveFilePath(str)
		assert.NotNil(err)
	}
}
