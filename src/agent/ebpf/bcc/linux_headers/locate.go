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

package linux_headers

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/tar"
)

var (
	lastInit       bool  = false
	lastInitResult error = nil
)

// Locate returns the path of the directory contains all of the Linux Kernel headers.
func locateHeader(libModuleDir string) (string, error) {
	// bcc/loader.cc looks for Linux headers in the following order:
	//   /lib/modules/<uname>/source
	//   /lib/modules/<uname>/build

	sourceDir := path.Join(libModuleDir, "source")
	buildDir := path.Join(libModuleDir, "build")

	if file.Exists(sourceDir) {
		return sourceDir, nil
	} else if file.Exists(buildDir) {
		return buildDir, nil
	}
	return "", fmt.Errorf("could not find 'source' or 'build' under %s", libModuleDir)
}

// locateClosestHeader returns the path of the directory contains all of the closest version Linux Kernel headers.
func locateClosestHeader(packageHeaderDir string, ver Version) (*Version, string, error) {
	files := file.List(packageHeaderDir)
	if files == nil {
		return nil, "", fmt.Errorf("could not find any file under %s", packageHeaderDir)
	}

	// packaged linux header dir has following struct
	// $package_dir/linux-header-<uname>.tar.gz

	var closestVersion *Version
	closestVersionPath := ""

	for _, file := range files {
		version, err := getKernelVersionFromArchiveFilePath(file)
		if err != nil {
			continue
		}

		if closestVersion == nil {
			closestVersion = &version
			closestVersionPath = file
			continue
		}

		if distance(ver, version) < distance(ver, *closestVersion) {
			closestVersion = &version
			closestVersionPath = file
		}
	}

	if closestVersion == nil {
		return nil, "", fmt.Errorf("could not find any kernel header tar file under %s", packageHeaderDir)
	}

	closestVersionFullPath := path.Join(packageHeaderDir, closestVersionPath)
	return closestVersion, closestVersionFullPath, nil
}

// installPackedHeader install kernel header from tar.gz file
func installPackagedHeader(packagedHeaderPath, installPath string) error {
	if !file.Exists(packagedHeaderPath) {
		return fmt.Errorf("packaged header %s not exist", packagedHeaderPath)
	}

	err := tar.GZExtract(packagedHeaderPath, installPath)
	if err != nil {
		return fmt.Errorf("decompress %s failed err: %v", packagedHeaderPath, err)
	}

	return nil
}

// modifyKernelVersion modify kernel version in package header file
func modifyKernelVersion(packageHeaderDir string, version Version) error {
	versionFilePath := path.Join(packageHeaderDir, "include/generated/uapi/linux/version.h")
	if !file.Exists(versionFilePath) {
		return fmt.Errorf("while modify kernel version, package header %s not exist", versionFilePath)
	}

	return WriteVersion(versionFilePath, version)
}

// findKernelConfig find kernel config file in the floowing order
// 1. /<host_root>/proc/config
// 2. /<host_root>/proc/config.gz
// 3. /<host_root>/boot/config-<uname>
// 4. /<host_root>/lib/modules/<uname>/config
// 5. /proc/config
// 6. /proc/config.gz
// return kernel config file path
func findKernelConfig(hostRootDir string, version Version, unameStr string) (string, error) {
	// Used when CONFIG_IKCONFIG=y is set.
	configPath := path.Join(hostRootDir, "proc/config")
	if file.Exists(configPath) {
		return configPath, nil
	}

	// Used when CONFIG_IKCONFIG_PROC=y is set.
	configPath = path.Join(hostRootDir, "proc/config.gz")
	if file.Exists(configPath) {
		return configPath, nil
	}

	// /boot/config-<uname>: Common place to store the config.
	configPath = path.Join(hostRootDir, fmt.Sprintf("/boot/config-%s", unameStr))
	if file.Exists(configPath) {
		return configPath, nil
	}

	// /lib/modules/<uname>/config: Used by RHEL8 CoreOS,
	configPath = path.Join(hostRootDir, fmt.Sprintf("lib/modules/%s/config", unameStr))
	if file.Exists(configPath) {
		return configPath, nil
	}

	if file.Exists("/proc/config") {
		return "/proc/config", nil
	}

	if file.Exists("/proc/config.gz") {
		return "/proc/config.gz", nil
	}

	return "", fmt.Errorf("no kernel config found")
}

// genAutoConf generate auto conf base on kernel config
// kernel config:
//  CONFIG_CC_VERSION_TEXT="gcc (GCC) 12.2.0"
//  CONFIG_CC_IS_GCC=y
//  CONFIG_GCC_VERSION=120200
// autoconf.h
//  #define CONFIG_CC_VERSION_TEXT "gcc (GCC) 12.2.0"
//  #define CONFIG_CC_IS_GCC 1
//  #define CONFIG_GCC_VERSION 120200
func genAutoConf(packageHeaderDir, configFilePath string) (int, error) {
	if !file.Exists(packageHeaderDir) {
		return 0, fmt.Errorf("empty package header dir %s", packageHeaderDir)
	}
	autoConfPath := path.Join(packageHeaderDir, "include/generated/autoconf.h")
	reader, closer, err := file.Reader(configFilePath)
	if err != nil {
		return 0, fmt.Errorf("could not reader %s error: %v", configFilePath, err)
	}

	defer closer.Close()

	if strings.HasSuffix(configFilePath, ".gz") {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return 0, fmt.Errorf("could not craete gzip reader %s error: %v", configFilePath, err)
		}
	}

	var hz int
	var lines []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		content := scanner.Text()
		if strings.HasPrefix(content, "CONFIG_HZ=") {
			hz, err = strconv.Atoi(strings.TrimPrefix(content, "CONFIG_HZ="))
			if err != nil {
				return 0, fmt.Errorf("could not parse reader %s CONFIG_HZ line error: %v", configFilePath, err)
			}
		}
		if strings.HasPrefix(content, "CONFIG_") {
			content = strings.ReplaceAll(content, "=y", " 1")
			content = strings.ReplaceAll(content, "=m", "_MODULE 1")
			content = strings.ReplaceAll(content, "=", " ")
			lines = append(lines, "#define "+content)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("could not reading reader %s error: %v", configFilePath, err)
	}

	contents := strings.Join(lines, "\n")

	err = file.Write(autoConfPath, contents)
	if err != nil {
		return 0, fmt.Errorf("could not write autoconf to %s error: %v", autoConfPath, err)
	}

	return hz, nil
}

// genTimeConst generate timeconst.h base on kernel config
func genTimeConst(hostRootDir, packageHeaderDir string, hz int) error {
	timeConstPath := path.Join(packageHeaderDir, "include/generated/timeconst.h")

	srcConstPath := hostRootDir + "/timeconst_" + strconv.Itoa(hz) + ".h"

	err := file.Copy(srcConstPath, timeConstPath)
	if err != nil {
		return fmt.Errorf("could not copy timeconst %s to %s error: %v", srcConstPath, timeConstPath, err)
	}
	return nil
}

// applyConfigPatches apply config patches to package header
func applyConfigPatches(hostRootDir, packageHeaderDir, starShipDir, unameStr string, version Version) error {
	kernelConfig, err := findKernelConfig(hostRootDir, version, unameStr)
	if err != nil {
		return err
	}

	hz, err := genAutoConf(packageHeaderDir, kernelConfig)
	if err != nil {
		return err
	}

	return genTimeConst(starShipDir, packageHeaderDir, hz)
}

// locateAndInstallPackagedHeaders that will find the Linux Kernel headers in the following order:
// 1. search closest version from packaged header directory
// 2. extract it to "/usr/src/linux-headers-<version>-starship" directory
// 3. modify kernel version in "/usr/src/linux-headers-<version>-starship/include/generated/uapi/linux/version.h"
// 4. apply config patches in "/usr/src/linux-headers-<version>-starship/include/generated/autoconf.h"
//    and "/usr/src/linux-headers-<version>-starship/include/generated/timeconst.h"
// 5. create a symlink from "/usr/src/linux-headers-<version>-starship" to "/usr/src/linux-headers-<version>/build"
func locateAndInstallPackageHeaders(hostRootDir, libModuleDir, starShipDir,
	installHeadersDir, unameStr string, version Version,
) error {
	libModuleBuildDir := path.Join(libModuleDir, "build")

	chosenVersion, chosenVersionPath, err := locateClosestHeader(starShipDir, version)
	if err != nil {
		return fmt.Errorf("while installing linux header, failed to locate the closest headers archive, error: %v", err)
	}

	installHeaderSubDir := fmt.Sprintf("usr/src/linux-headers-%s-starship", chosenVersion.semVerStr())
	installHeaderDir := path.Join(installHeadersDir, installHeaderSubDir)

	err = installPackagedHeader(chosenVersionPath, installHeadersDir)
	if err != nil {
		return fmt.Errorf("while installing linux header, failed to install oackaged headers, error: %v", err)
	}

	err = modifyKernelVersion(installHeaderDir, version)
	if err != nil {
		return fmt.Errorf("while installing linux header, failed to modify kernel version, error: %v", err)
	}

	err = applyConfigPatches(hostRootDir, installHeaderDir, starShipDir, unameStr, version)
	if err != nil {
		return fmt.Errorf("while installing linux header, failed to apply config, error: %v", err)
	}

	if err = file.CreateSymLink(installHeaderDir, libModuleBuildDir); err != nil {
		return fmt.Errorf("while installing linux header, failed to create symlink, error: %v", err)
	}

	return nil
}

// locateAndLinkHostHeader search host mount dir /host/lib/modules/<uname -r>/
// and symlink to /lib/modules/<uname -r>/build
// hostRootDir: /host is the host mount dir
// libModuleDir: /lib/modules/<uname -r> is the module dir in container
// hostModuleDir: /host/lib/modules/<uname -r> is the module dir in host
func locateAndLinkHostHeader(hostRootDir, libModuleDir, hostModuleDir string) error {
	//  searching /host/lib/modules/<uname>/source and try to link to /lib/modules/<uname>/source
	hostHeaderSourceDir := path.Join(hostModuleDir, "source")
	linkRealDir, err := file.ReadSymLink(hostHeaderSourceDir)
	if err == nil {
		hostlinkRealDir := path.Join(hostRootDir, linkRealDir)
		libModuleSourceDir := path.Join(libModuleDir, "source")
		if err := file.CreateSymLink(hostlinkRealDir, libModuleSourceDir); err != nil {
			return fmt.Errorf(
				"while locate and install host linux header, failed to create symlink for source, error: %v", err)
		}
		return nil
	}

	// searching /host/lib/modules/<uname>/build and try to link to /lib/modules/<uname>/build
	hostHeaderBuildDir := path.Join(hostModuleDir, "build")
	linkRealDir, err = file.ReadSymLink(hostHeaderBuildDir)
	if err == nil {
		hostlinkRealDir := path.Join(hostRootDir, linkRealDir)
		libModuleBuildDir := path.Join(libModuleDir, "build")
		if err := file.CreateSymLink(hostlinkRealDir, libModuleBuildDir); err != nil {
			return fmt.Errorf(
				"while locate and install host linux header,failed to create symlink for build , error: %v", err)
		}
		return nil
	}

	return fmt.Errorf(
		"while locate and install host linux header, could not find 'source' or 'build' under %s", hostModuleDir)
}

// findOrInstallHeaders that will find the Linux Kernel headers in the following order:
// 1. search locale container linux header in '/lib/modules/<uname -r>/' dir
// 2. search host mount dir /host/lib/modules/<uname -r>/, install host linux header
// 3. search starship packaged header, install starship linux header
func findOrInstallHeaders(hostRootDir, starShipDir, libModulesDir, installHeadersDir string) error {
	version, err := GetVersion()
	if err != nil {
		return err
	}

	uStr, err := unameStr()
	if err != nil {
		return err
	}

	libModuleDir := path.Join(libModulesDir, uStr)
	if _, err = locateHeader(libModuleDir); err == nil {
		return nil
	}

	hostModuleDir := path.Join(hostRootDir, "lib/modules", uStr)
	if err := locateAndLinkHostHeader(hostRootDir, libModuleDir, hostModuleDir); err == nil {
		return nil
	}

	return locateAndInstallPackageHeaders(hostRootDir, libModuleDir, starShipDir, installHeadersDir, uStr, version)
}

// Init will find or install linux headers
func Init() error {
	if lastInit {
		return lastInitResult
	}

	// hostRootDir refers to the path to the mounted volume which points the host's root path /.
	// starShipDir is the starship header compressed tar dir
	// libModulesDir is the module dir in container, BCC will search the header in this dir
	// installHeadersDir is the install dir in container, because of the
	// compressed tar struct is /usr/src/linux-headers-<version>-starship, so compress destination is /
	hostRootDir := "/host"
	starShipDir := "/starship/linux_headers"
	libModulesDir := "/lib/modules"
	installHeadersDir := "/"
	err := findOrInstallHeaders(hostRootDir, starShipDir, libModulesDir, installHeadersDir)

	lastInit = true

	if err != nil {
		lastInitResult = fmt.Errorf("failed to find or install linux headers, error: %v", err)
		return lastInitResult
	}

	return nil
}
