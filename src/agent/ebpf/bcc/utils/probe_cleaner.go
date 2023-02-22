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

package utils

import (
	"fmt"
	"path"
	"strings"

	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/log"
)

const (
	// The path of kprobe files under /sys, join with the host sys root path to form the correct path inside container.
	kprobeEventsSysRelPath = "kernel/debug/tracing/kprobe_events"
	uprobeEventsSysRelPath = "kernel/debug/tracing/uprobe_events"

	// Marker is set at:
	// https://github.com/Tricorder Observability/bcc/commit/50de7107d6a48fcfe4f82d33433960f965d1a16a
	//
	// This ensures we only removing Tricorder attached probes, and not affecting the probes attached by other users
	// using different tools.
	//
	// This only applies to BCC with Debugfs (an alternative is BTF? TODO(yzhao): Needs more investigation):
	// https://github.com/iovisor/bcc/blob/master/INSTALL.md#setup-required-to-run-the-tools
	tricorderMarker = "__tricorder__"
)

// findProbes returns the lines of probes in probeFile which contains marker.
func findProbes(probeFile string, marker string) ([]string, error) {
	fileContent, err := file.Read(probeFile)
	if err != nil {
		return nil, fmt.Errorf("while searching for attached probes with marker '%s', failed to read file '%s', error: %v",
			marker, probeFile, err)
	}
	res := make([]string, 0)
	for _, line := range strings.Split(fileContent, "\n") {
		if !strings.Contains(line, marker) {
			continue
		}
		fields := strings.Split(line, " ")
		if len(fields) != 2 {
			continue
		}
		// Note that a probe looks like the following:
		//     [p|r]:kprobes/your_favorite_probe_name_here __x64_sys_connect
		// https://www.kernel.org/doc/html/latest/trace/kprobetrace.html
		probeName := fields[0]
		if probeName[0] != 'p' && probeName[0] != 'r' {
			continue
		}
		res = append(res, probeName)
	}
	return res, nil
}

// cleanProbes writes lines that would remove the input probes.
func cleanProbes(probeFile string, probes []string) error {
	lines := make([]string, 0, len(probes))
	for _, probe := range probes {
		parts := strings.SplitN(probe, ":", 2)
		if len(parts) != 2 {
			continue
		}
		deleteProbeText := "-:" + parts[1]
		lines = append(lines, deleteProbeText)
	}
	err := file.Append(probeFile, strings.Join(lines, "\n"))
	if err != nil {
		return fmt.Errorf("while cleaning probes, failed to write delete probe text, error: %v", err)
	}
	return nil
}

func findAndCleanProbes(probeFile string, marker string) error {
	probes, err := findProbes(probeFile, tricorderMarker)
	if err != nil {
		return fmt.Errorf("while cleaning probes, failed to find relevant kprobes, error: %v", err)
	}

	origCount := len(probes)

	err = cleanProbes(probeFile, probes)
	if err != nil {
		return fmt.Errorf("while cleaning probes, failed to clean kprobes, error: %v", err)
	}

	probes, err = findProbes(probeFile, tricorderMarker)
	afterCleanCount := len(probes)
	if err == nil {
		log.Infof("Found %d probes in %s with marker %s, %d probes left after cleaning",
			origCount, probeFile, marker, afterCleanCount)
	}

	return nil
}

// CleanTricorderProbes removes all of the probes attached by tricorder.
func CleanTricorderProbes(hostSysRootPath string) error {
	if err := findAndCleanProbes(path.Join(hostSysRootPath, kprobeEventsSysRelPath), tricorderMarker); err != nil {
		return err
	}
	if err := findAndCleanProbes(path.Join(hostSysRootPath, uprobeEventsSysRelPath), tricorderMarker); err != nil {
		return err
	}
	return nil
}
