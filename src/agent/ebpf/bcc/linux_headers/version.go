// Copyright (C) 2023 Tricorder Observability
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

package linux_headers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/tricorder/src/utils/common"
	"github.com/tricorder/src/utils/file"
)

// Version represents the system's kernel version.
type Version struct {
	ver   uint8
	major uint8
	minor uint16
}

// This code produces the same result as what's defined in the Linux version.h header
// #define KERNEL_VERSION(a,b,c) (((a) << 16) + ((b) << 8) + (c))
func (v Version) code() uint32 {
	return uint32(v.ver)<<16 + uint32(v.major)<<8 + uint32(v.minor)
}

// semVerStr returns a semantic version string.
func (v Version) semVerStr() string {
	return fmt.Sprintf("%d.%d.%d", v.ver, v.major, v.minor)
}

// distance returns a numeric value indicating the differences between 2 versions.
func distance(v1, v2 Version) int {
	return common.AbsUint16s(
		v1.minor, v2.minor) + 100*common.AbsUint8s(v1.major, v2.major) + 10000*common.AbsUint8s(v1.ver, v2.ver)
}

func parseVersion(str string) (Version, error) {
	// Regexp grammar is in https://pkg.go.dev/regexp/syntax
	r := regexp.MustCompile(`^(?P<ver>[0-9]+)\.(?P<major>[0-9]+)\.(?P<minor>[0-9]+)`)
	matches := r.FindStringSubmatch(str)
	// First element is the left-most match of entire regexp, 1:4 are the match groups.
	if len(matches) != 4 {
		return Version{}, fmt.Errorf("while parsing version string, failed to get 3 dot-separated fields from '%s'", str)
	}
	ver, err := strconv.Atoi(matches[1])
	if err != nil {
		return Version{}, fmt.Errorf("while parsing version string, 1st component is not number '%s'", str)
	}
	major, err := strconv.Atoi(matches[2])
	if err != nil {
		return Version{}, fmt.Errorf("while parsing version string, 2nd component is not number '%s'", str)
	}
	minor, err := strconv.Atoi(matches[3])
	if err != nil {
		return Version{}, fmt.Errorf("while parsing version string, 3rd component is not number '%s'", str)
	}
	return Version{
		ver:   uint8(ver),
		major: uint8(major),
		minor: uint16(minor),
	}, nil
}

// utsnameToString converts the utsname to a string and returns it.
// Since the Utsname fields are all fixed length byte array.
func utsnameToString(unameArray [65]byte) string {
	l := len(unameArray)
	// Find the null char in the byte array.
	for i := 0; i < len(unameArray); i++ {
		if unameArray[i] == 0 {
			l = i
			break
		}
	}
	return string(unameArray[:l])
}

// unameStr return system version string, effectively running `uname -r`.
func unameStr() (string, error) {
	uname := new(unix.Utsname)
	err := unix.Uname(uname)
	if err != nil {
		return "", fmt.Errorf("while getting system version, failed to call uname, error: %v", err)
	}
	return utsnameToString(uname.Release), nil
}

// unameVersion return system version release, effectively running `uname -r`.
func unameVersion() (Version, error) {
	uname, err := unameStr()
	if err != nil {
		return Version{}, err
	}
	return parseVersion(uname)
}

// procVersion returns system version from /proc/version_signature
func procVersion() (Version, error) {
	const path = "/proc/version_signature"
	content, err := file.Read(path)
	if err != nil {
		return Version{}, fmt.Errorf("while getting version from /proc, failed to read '%s', error: %v", path, err)
	}
	tokens := strings.Split(content, " ")
	return parseVersion(tokens[len(tokens)-1])
}

func GetVersion() (Version, error) {
	ver, err := unameVersion()
	if err == nil {
		return ver, nil
	}
	ver, err = procVersion()
	if err == nil {
		return ver, nil
	}
	return Version{}, fmt.Errorf("failed to get kernel version")
}

// WriteVersion replaces the macro in versionHeaderFile such that the desired version is as the specified version
// by the input version.
func WriteVersion(versionHeaderFile string, ver Version) error {
	const macroPrefix = "#define LINUX_VERSION_CODE "
	desiredVerStr := macroPrefix + strconv.Itoa(int(ver.code()))

	content, err := file.Read(versionHeaderFile)
	if err != nil {
		return fmt.Errorf("while writing version to '%s', failed to read its content, error: %v", versionHeaderFile, err)
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, macroPrefix) {
			lines[i] = desiredVerStr
		}
	}
	return file.Write(versionHeaderFile, strings.Join(lines, "\n"))
}

// getKernelVersionFromArchiveFilePath return Version from archive file path
func getKernelVersionFromArchiveFilePath(file string) (Version, error) {
	if !strings.HasPrefix(file, "linux-headers-") {
		return Version{}, fmt.Errorf("no 'linux-headers-' prefix, failed to parse version from file %s", file)
	}
	if !strings.HasSuffix(file, ".tar.gz") {
		return Version{}, fmt.Errorf("no '.tar.gz' suffix, failed to parse version from file %s", file)
	}

	versionStr := common.StrTrimPrefix(file, len("linux-headers-"))
	versionStr = common.StrTrimSuffix(versionStr, len(".tar.gz"))
	return parseVersion(versionStr)
}
