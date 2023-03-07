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

package linux_headers

import (
	"fmt"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	testuitls "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/file"
)

// Tests that locateHeader() returns the correct header dir.
func TestLocateHeader(t *testing.T) {
	assert := assert.New(t)

	// module dir 'is non-existent-dir', return error
	assert.NotNil(locateHeader("non-existent-dir"))

	// module dir just contains a 'source' dir, return 'source' dir
	tmpDir := testuitls.CreateTmpDir()
	tmpSourceDir := path.Join(tmpDir, "source")
	assert.Nil(file.CreateDir(tmpSourceDir))
	resDir, err := locateHeader(tmpDir)
	assert.Equal(tmpSourceDir, resDir)
	assert.Nil(err)

	// module dir just contains a 'build' dir, return 'build' dir
	tmpDir = testuitls.CreateTmpDir()
	tmpBuildDir := path.Join(tmpDir, "build")
	assert.Nil(file.CreateDir(tmpBuildDir))
	resDir, err = locateHeader(tmpDir)
	assert.Equal(tmpBuildDir, resDir)
	assert.Nil(err)

	// module dir contains both 'source' and 'build' dir, return 'source' dir
	tmpDir = testuitls.CreateTmpDir()
	tmpSourceDir = path.Join(tmpDir, "source")
	tmpBuildDir = path.Join(tmpDir, "build")
	assert.Nil(file.CreateDir(tmpBuildDir))
	assert.Nil(file.CreateDir(tmpSourceDir))
	resDir, err = locateHeader(tmpDir)
	assert.Equal(tmpSourceDir, resDir)
	assert.Nil(err)

	// module dir does not contain 'source' or 'build' dir, return error
	tmpDir = testuitls.CreateTmpDir()
	resDir, err = locateHeader(tmpDir)
	assert.Equal("", resDir)
	assert.ErrorContains(err, "could not find 'source' or 'build' under ", tmpDir)
}

// Tests that locateClosestHeader() returns the correct packed header path.
func TestLocateClosestHeader(t *testing.T) {
	assert := assert.New(t)

	// packaged header dir 'is non-existent-dir', return error
	resVer, resPath, err := locateClosestHeader("non-existent-dir", Version{})
	assert.Nil(resVer)
	assert.Equal("", resPath)
	assert.NotNil(err)

	// paclaged header dir is empty, return error
	tmpDir := testuitls.CreateTmpDir()
	resVer, resPath, err = locateClosestHeader(tmpDir, Version{})
	assert.Nil(resVer)
	assert.Equal("", resPath)
	assert.ErrorContains(err, "could not find any kernel header tar file under", tmpDir)

	// packaged header dir contains wrong name header tar file, return error
	tmpDir = testuitls.CreateTmpDir()
	tmpTarPath := path.Join(tmpDir, "test.tar.gz")
	assert.Nil(file.Create(tmpTarPath))
	resVer, resPath, err = locateClosestHeader(tmpDir, Version{})
	assert.Nil(resVer)
	assert.Equal("", resPath)
	assert.ErrorContains(err, "could not find any kernel header tar file under", tmpDir)

	// packaged header dir contains correct one name header tar file, return the tar file
	tmpDir = testuitls.CreateTmpDir()
	tmpTarPath = path.Join(tmpDir, "linux-headers-4.4.0-116.tar.gz")
	assert.Nil(file.Create(tmpTarPath))
	resVer, resPath, err = locateClosestHeader(tmpDir, Version{})
	assert.Equal(Version{4, 4, 0}, *resVer)
	assert.Equal(tmpTarPath, resPath)
	assert.Nil(err)

	// packaged header dir contains correct two name header tar file, return the closest one
	tmpDir = testuitls.CreateTmpDir()
	tmpTarPathVersion4 := path.Join(tmpDir, "linux-headers-4.4.0-116.tar.gz")
	assert.Nil(file.Create(tmpTarPathVersion4))
	tmpTarPathVersion5 := path.Join(tmpDir, "linux-headers-5.5.0-117.tar.gz")
	assert.Nil(file.Create(tmpTarPathVersion5))

	resVer, resPath, err = locateClosestHeader(tmpDir, Version{4, 5, 0})
	assert.Equal(Version{4, 4, 0}, *resVer)
	assert.Equal(tmpTarPathVersion4, resPath)
	assert.Nil(err)

	resVer, resPath, err = locateClosestHeader(tmpDir, Version{5, 6, 0})
	assert.Equal(Version{5, 5, 0}, *resVer)
	assert.Equal(tmpTarPathVersion5, resPath)
	assert.Nil(err)

	resVer, resPath, err = locateClosestHeader(tmpDir, Version{6, 6, 0})
	assert.Equal(Version{5, 5, 0}, *resVer)
	assert.Equal(tmpTarPathVersion5, resPath)
	assert.Nil(err)
}

// Tests that installPackagedHeader() installs the correct header.
func TestInstallPackagedHeader(t *testing.T) {
	assert := assert.New(t)

	// packaged header is 'is non-existent-file', return error
	assert.NotNil(installPackagedHeader("non-existent-file", "tmp"))

	// packaged header is wrong type file , return error
	tmpDir := testuitls.CreateTmpDir()
	tmpFile := testuitls.TestFilePath("src/utils/tar/testdata/wrong_file_format.tar.gz")

	err := installPackagedHeader(tmpFile, tmpDir)
	assert.ErrorContains(err, "decompress ", tmpFile, " failed err:")

	// packaged header is correct file, install the header to the correct dir
	tmpDir = testuitls.CreateTmpDir()
	tmpFile = testuitls.TestFilePath("src/utils/tar/testdata/test.tar.gz")
	err = installPackagedHeader(tmpFile, tmpDir)
	assert.Nil(err)
	assert.True(file.Exists(path.Join(tmpDir, "hello.txt")))
}

// Tests that modifyKernelVersion() modifies the correct 'include/generated/uapi/linux/version.h' file.
func TestModifyKernelVersion(t *testing.T) {
	assert := assert.New(t)

	// packaged header is 'non-existent-dir', return error
	assert.NotNil(modifyKernelVersion("non-existent-dir", Version{}))

	// 'include/generated/uapi/linux/version.h' is correct file, modify the file
	tmpFile := testuitls.TestFilePath("devops/linux_headers/output/linux-headers-5.2.1-starship.tar.gz")
	tmpDir := testuitls.CreateTmpDir()
	assert.Nil(installPackagedHeader(tmpFile, tmpDir))
	targetVersion := Version{4, 4, 0}
	tmpHeadersDir := path.Join(tmpDir, "usr/src/linux-headers-5.2.1-starship")
	assert.Nil(modifyKernelVersion(tmpHeadersDir, targetVersion))

	writeFile := path.Join(tmpHeadersDir, "include/generated/uapi/linux/version.h")
	assert.True(file.Exists(writeFile))
	assert.True(file.Contains(writeFile, fmt.Sprintf("#define LINUX_VERSION_CODE %s",
		strconv.Itoa(int(targetVersion.code())))))
}

// Tests that findKernelConfig() returns the correct kernel config file path.
func TestFindKernelConfig(t *testing.T) {
	assert := assert.New(t)

	// find kernel config file in 'non-existent-dir', return error
	assert.NotNil(findKernelConfig("non-existent-dir", Version{}, "5.15.0-58-generic"))

	// kernel config in /<host_root>/proc/config is correct, return the correct path
	tmpDir := testuitls.CreateTmpDir()
	tmpConfigFile := path.Join(tmpDir, "proc/config")
	assert.Nil(file.Create(tmpConfigFile))
	resPath, err := findKernelConfig(tmpDir, Version{}, "5.15.0-58-generic")
	assert.Equal(tmpConfigFile, resPath)
	assert.Nil(err)

	// kernel config in /<host_root>/proc/config.gz is correct, return the correct path
	tmpDir = testuitls.CreateTmpDir()
	tmpConfigFile = path.Join(tmpDir, "proc/config.gz")
	assert.Nil(file.Create(tmpConfigFile))
	resPath, err = findKernelConfig(tmpDir, Version{}, "5.15.0-58-generic")
	assert.Equal(tmpConfigFile, resPath)
	assert.Nil(err)

	// kernel config in /<host_root>/proc/config and /<host_root>/proc/config.gz are both correct,
	// return the correct path
	tmpDir = testuitls.CreateTmpDir()
	tmpConfigFile = path.Join(tmpDir, "proc/config")
	tmpConfigGZFile := path.Join(tmpDir, "proc/config.gz")
	assert.Nil(file.Create(tmpConfigFile))
	assert.Nil(file.Create(tmpConfigGZFile))
	resPath, err = findKernelConfig(tmpDir, Version{}, "5.15.0-58-generic")
	assert.Equal(tmpConfigFile, resPath)
	assert.Nil(err)
}

// Tests that genAutoConf() generates the correct autoconf file.
func TestGenAutoConf(t *testing.T) {
	assert := assert.New(t)

	tmpDir := testuitls.CreateTmpDir()
	// config file path is 'non-existent-dir', return error
	resHZ, err := genAutoConf(tmpDir, "non-existent-dir")
	assert.Equal(0, resHZ)
	assert.ErrorContains(err, "could not reader ", path.Join("non-existent-dir", "include/generated/autoconf.h"))

	// package header is 'non-existent-dir', return error
	configPath := testuitls.TestFilePath("src/agent/ebpf/bcc/linux_headers/testdata/config")
	resHZ, err = genAutoConf("non-existent-dir", configPath)
	assert.Equal(0, resHZ)
	assert.ErrorContains(err, "empty package header dir ", "non-existent-dir")

	// config file is correct, generate the correct autoconf file
	tmpDir = testuitls.CreateTmpDir()
	configPath = testuitls.TestFilePath("src/agent/ebpf/bcc/linux_headers/testdata/config")
	resHZ, err = genAutoConf(tmpDir, configPath)
	assert.Nil(err)
	assert.Equal(250, resHZ)
	assert.Nil(err)
	content, err := file.Read(path.Join(tmpDir, "include/generated/autoconf.h"))
	assert.Nil(err)
	expectedString := `#define CONFIG_CRYPTO_HMAC 1
#define CONFIG_CRYPTO_XCBC_MODULE 1
#define CONFIG_CRYPTO_VMAC_MODULE 1
#define CONFIG_CRYPTO_CRC32C 1
#define CONFIG_CRYPTO_CRC32C_INTEL 1
#define CONFIG_CRYPTO_CRC32_MODULE 1`
	assert.Contains(content, expectedString)

	// config.gz file is correct, generate the correct autoconf file
	tmpDir = testuitls.CreateTmpDir()
	configPath = testuitls.TestFilePath("src/agent/ebpf/bcc/linux_headers/testdata/config.gz")
	resHZ, err = genAutoConf(tmpDir, configPath)
	assert.Nil(err)
	assert.Equal(250, resHZ)
	assert.Nil(err)
	content, err = file.Read(path.Join(tmpDir, "include/generated/autoconf.h"))
	assert.Nil(err)
	assert.Contains(content, expectedString)
}

// Tests that genTimeConst() generates the correct time constant file.
func TestGenTimeConst(t *testing.T) {
	assert := assert.New(t)

	tmpDir := testuitls.CreateTmpDir()
	// host mount dir is 'non-existent-dir', return error
	assert.NotNil(genTimeConst("non-existent-dir", tmpDir, 250))

	// package header is 'non-existent-dir', return error
	assert.NotNil(genTimeConst(tmpDir, "non-existent-dir", 250))

	// HZ is 0, return error
	tmpHostDir := testuitls.CreateTmpDir()
	tmpPackageDir := testuitls.CreateTmpDir()
	timeConstPath := path.Join(tmpHostDir, "timeconst_250.h")
	assert.Nil(file.Create(timeConstPath))
	assert.NotNil(genTimeConst(tmpHostDir, tmpPackageDir, 0))

	// HZ is 250, return correct time constant file
	tmpHostDir = testuitls.CreateTmpDir()
	tmpPackageDir = testuitls.CreateTmpDir()
	timeConstPath = path.Join(tmpHostDir, "timeconst_250.h")
	tmpTimeConstPath := testuitls.TestFilePath("devops/linux_headers/output/timeconst_250.h")
	assert.Nil(file.Copy(tmpTimeConstPath, timeConstPath))
	assert.Nil(genTimeConst(tmpHostDir, tmpPackageDir, 250))
	content, err := file.Read(path.Join(tmpTimeConstPath))
	assert.Nil(err)
	tmpContent, err := file.Read(path.Join(tmpPackageDir, "include/generated/timeconst.h"))
	assert.Nil(err)
	assert.Equal(content, tmpContent)
}

// Tests that applyConfigPatches() applies the correct kernel config.
// The test is based on the kernel config file in the testdata directory.
func TestApplyConfigPatches(t *testing.T) {
	assert := assert.New(t)

	tmpHostRootDir := testuitls.CreateTmpDir()
	tmpPackageHeaderDir := testuitls.CreateTmpDir()
	tmpStarShipDir := testuitls.CreateTmpDir()

	// config file path is 'non-existent-dir', return error
	err := applyConfigPatches("non-existent-dir", tmpPackageHeaderDir,
		tmpStarShipDir, "5.4.0-58-generic", Version{5, 4, 0})
	assert.ErrorContains(err, "no kernel config found")

	configPath := path.Join(tmpHostRootDir, "proc/config")
	testConfigPath := testuitls.TestFilePath("src/agent/ebpf/bcc/linux_headers/testdata/config")
	assert.Nil(file.Create(configPath))
	assert.Nil(file.Copy(testConfigPath, configPath))

	// package header path is 'non-existent-dir', return error
	err = applyConfigPatches(tmpHostRootDir, "non-existent-dir",
		tmpStarShipDir, "5.4.0-58-generic", Version{5, 4, 0})
	assert.ErrorContains(err, "empty package header dir")

	// starship path is 'non-existent-dir', return error
	err = applyConfigPatches(tmpHostRootDir, tmpPackageHeaderDir,
		"non-existent-dir", "5.4.0-58-generic", Version{5, 4, 0})
	assert.ErrorContains(err, "could not copy timeconst")

	tmpPackageHeaderDir = testuitls.CreateTmpDir()
	timeConstPath := path.Join(tmpStarShipDir, "timeconst_250.h")
	testTimeConstPath := testuitls.TestFilePath("devops/linux_headers/output/timeconst_250.h")
	assert.Nil(file.Create(timeConstPath))
	assert.Nil(file.Copy(testTimeConstPath, timeConstPath))

	// all dependencies are correct, return nil
	assert.Nil(applyConfigPatches(tmpHostRootDir, tmpPackageHeaderDir,
		tmpStarShipDir, "5.4.0-58-generic", Version{5, 4, 0}))
	assert.Equal(true, file.Exists(path.Join(tmpPackageHeaderDir, "include/generated/timeconst.h")))
	assert.Equal(true, file.Exists(path.Join(tmpPackageHeaderDir, "include/generated/autoconf.h")))
}

// Test that locateAndInstallPackageHeaders() locates and installs the correct package headers.
func TestLocateAndInstallPackageHeaders(t *testing.T) {
	assert := assert.New(t)

	tmpHostRootDir := testuitls.CreateTmpDir()
	tmpLibModuleDir := testuitls.CreateTmpDir()
	tmpStarShipDir := testuitls.CreateTmpDir()
	tmpInstallHeadersDir := testuitls.CreateTmpDir()

	configPath := path.Join(tmpHostRootDir, "proc/config")
	testConfigPath := testuitls.TestFilePath("src/agent/ebpf/bcc/linux_headers/testdata/config")
	assert.Nil(file.Create(configPath))
	assert.Nil(file.Copy(testConfigPath, configPath))

	timeConstPath := path.Join(tmpStarShipDir, "timeconst_250.h")
	testTimeConstPath := testuitls.TestFilePath("devops/linux_headers/output/timeconst_250.h")
	assert.Nil(file.Create(timeConstPath))
	assert.Nil(file.Copy(testTimeConstPath, timeConstPath))

	headersGZ := path.Join(tmpStarShipDir, "linux-headers-5.1.1.tar.gz")
	testHeadersGZ := testuitls.TestFilePath("devops/linux_headers/output/linux-headers-5.1.1-starship.tar.gz")
	assert.Nil(file.Create(headersGZ))
	assert.Nil(file.Copy(testHeadersGZ, headersGZ))

	// host root dir is 'non-existent-dir', return error
	err := locateAndInstallPackageHeaders("non-existent-dir", tmpLibModuleDir, tmpStarShipDir,
		tmpInstallHeadersDir, "5.4.0-58-generic", Version{5, 4, 0})
	assert.ErrorContains(err, "no kernel config found")

	// lib module dir is 'non-existent-dir', no return error
	err = locateAndInstallPackageHeaders(tmpHostRootDir, "non-existent-dir", tmpStarShipDir,
		tmpInstallHeadersDir, "5.4.0-58-generic", Version{5, 4, 0})
	assert.Nil(err)
	assert.Equal(true, file.Exists("non-existent-dir/build"))

	// starship dir is 'non-existent-dir', return error
	err = locateAndInstallPackageHeaders(tmpHostRootDir, tmpLibModuleDir, "non-existent-dir",
		tmpInstallHeadersDir, "5.4.0-58-generic", Version{5, 4, 0})
	assert.ErrorContains(err, "could not find any kernel header tar file ")

	// install headers dir is 'non-existent-dir', no return error
	err = locateAndInstallPackageHeaders(tmpHostRootDir, tmpLibModuleDir, tmpStarShipDir,
		"non-existent-dir", "5.4.0-58-generic", Version{5, 4, 0})
	assert.Nil(err)
	assert.Equal(true, file.Exists("non-existent-dir/usr/src/linux-headers-5.1.1-starship"))

	// all dependencies are correct, return nil
	tmpLibModuleDir = testuitls.CreateTmpDir()
	tmpInstallHeadersDir = testuitls.CreateTmpDir()
	err = locateAndInstallPackageHeaders(tmpHostRootDir, tmpLibModuleDir, tmpStarShipDir,
		tmpInstallHeadersDir, "5.4.0-58-generic", Version{5, 4, 0})
	assert.Nil(err)
	assert.Equal(true, file.Exists(path.Join(tmpInstallHeadersDir, "usr/src/linux-headers-5.1.1-starship")))
	assert.Equal(true, file.Exists(path.Join(tmpLibModuleDir, "build/include/generated/timeconst.h")))
	assert.Equal(true, file.Exists(path.Join(tmpLibModuleDir, "build/include/generated/autoconf.h")))
}

// Tests that locateAndLinkHostHeader link the correct host header.
func TestLocateAndLinkHostHeader(t *testing.T) {
	assert := assert.New(t)

	tmpHostRootDir := testuitls.CreateTmpDir()
	tmpLibModuleDir := testuitls.CreateTmpDir()
	tmpRealModuleDir := testuitls.CreateTmpDir()

	tmpHostModuleDir := path.Join(tmpHostRootDir, "lib/modules/5.1.1-starship/")
	tmpHostModuleBuildDir := path.Join(tmpHostModuleDir, "build")
	assert.Nil(file.CreateSymLink(tmpRealModuleDir, tmpHostModuleBuildDir))

	// host root dir is 'non-existent-dir-1', return error
	err := locateAndLinkHostHeader("/", "non-existent-dir-1", tmpHostModuleDir)
	assert.Nil(err)

	// host module dir is 'non-existent-dir-2', return error
	err = locateAndLinkHostHeader("/", tmpLibModuleDir, "non-existent-dir-2")
	assert.ErrorContains(err, "could not find 'source' or 'build' under")

	// all dependencies are correct, return nil
	tmpHostRootDir = testuitls.CreateTmpDir()
	tmpLibModuleDir = testuitls.CreateTmpDir()
	tmpRealModuleDir = testuitls.CreateTmpDir()

	tmpHostModuleDir = path.Join(tmpHostRootDir, "lib/modules/5.1.1-starship")
	tmpHostModuleBuildDir = path.Join(tmpHostModuleDir, "build")
	assert.Nil(file.CreateSymLink(tmpRealModuleDir, tmpHostModuleBuildDir))
	err = locateAndLinkHostHeader("/", tmpLibModuleDir, tmpHostModuleDir)
	tmpLibModuleDir = path.Join(tmpLibModuleDir, "build")
	assert.Nil(err)
	assert.Equal(true, file.Exists(tmpLibModuleDir))
}

// Tests that findKernelHeadersTarFile() finds the correct kernel headers tar file.
func TestFindKernelHeadersTarFile(t *testing.T) {
	assert := assert.New(t)

	tmpHostRootDir := testuitls.CreateTmpDir()
	tmpLibModulesDir := testuitls.CreateTmpDir()
	tmpStarShipDir := testuitls.CreateTmpDir()
	tmpInstallHeadersDir := testuitls.CreateTmpDir()

	configPath := path.Join(tmpHostRootDir, "proc/config")
	testConfigPath := testuitls.TestFilePath("src/agent/ebpf/bcc/linux_headers/testdata/config")
	assert.Nil(file.Create(configPath))
	assert.Nil(file.Copy(testConfigPath, configPath))

	timeConstPath := path.Join(tmpStarShipDir, "timeconst_250.h")
	testTimeConstPath := testuitls.TestFilePath("devops/linux_headers/output/timeconst_250.h")
	assert.Nil(file.Create(timeConstPath))
	assert.Nil(file.Copy(testTimeConstPath, timeConstPath))

	headersGZ := path.Join(tmpStarShipDir, "linux-headers-5.1.1-starship.tar.gz")
	testHeadersGZ := testuitls.TestFilePath("devops/linux_headers/output/linux-headers-5.1.1-starship.tar.gz")
	assert.Nil(file.Create(headersGZ))
	assert.Nil(file.Copy(testHeadersGZ, headersGZ))

	uStr, err := unameStr()
	assert.Nil(err)

	assert.Nil(findOrInstallHeaders(tmpHostRootDir, tmpStarShipDir, tmpLibModulesDir, tmpInstallHeadersDir))
	assert.Equal(true, file.Exists(path.Join(tmpInstallHeadersDir, "usr/src/linux-headers-5.1.1-starship")))

	installLibModulesDir := path.Join(tmpLibModulesDir, uStr)
	assert.Equal(true, file.Exists(path.Join(installLibModulesDir, "build/include/generated/timeconst.h")))
}
