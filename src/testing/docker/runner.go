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

package testing

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/utils/common"
	"github.com/tricorder/src/utils/exec"
	"github.com/tricorder/src/utils/timer"
)

// Runner describes how to run an image.
type Runner struct {
	ContainerName string
	ImageName     string
	Options       []string
	Cmd           *exec.Command

	// A list of command line flags passed to application process
	// Note that this not the same as Options
	AppFlags []string

	// Look for this message in stdout and stderr before finishing launching.
	RdyMsg string
}

// Run starts the subprocess to run the declared command.
func (r *Runner) Launch(timeout time.Duration) error {
	cmdSlice := []string{
		"docker", "run",
		// This expose all declared port to any available host port.
		// This is later used to get the exposed port.
		"--publish-all",
		// Non-blocking execution, immediately return after launching.
		"--detach",
		// -it is a shorthand for --interactive --tty
		// Keep STDIN open even if not attached
		"--interactive",
		// pseudo-TTY
		"--tty",
		// Delete container when it's stopped.
		"--rm",
	}
	cmdSlice = append(cmdSlice, r.Options...)
	if len(r.ContainerName) == 0 {
		r.ContainerName = common.RandStr(24)
	}
	cmdSlice = append(cmdSlice, fmt.Sprintf("--name=%s", r.ContainerName))
	cmdSlice = append(cmdSlice, r.ImageName)
	cmdSlice = append(cmdSlice, r.AppFlags...)

	outStr, errStr, err := exec.Run(cmdSlice)
	log.Infof("command=%v, stdout=%s stderr=%s err=%v", cmdSlice, outStr, errStr, err)

	// timd.Tick() leaks the underlying Ticker, as it has no way to shutdown.
	// So this can only be used in test. Staticcheck is disabled to make this work in tests.

	timer := timer.New()

	// https://github.com/golangci/golangci-lint/issues/741
	//nolint:staticcheck // SA1015 this disable staticcheck for the next line
	c := time.Tick(time.Second)
	for range c {
		log.Infof("Looking for ready message: %s", r.RdyMsg)
		logs, err := r.Logs()
		if err != nil {
			return fmt.Errorf("while looking for ready message in the log, failed to get logs, error: %v", err)
		}
		if strings.Contains(logs, r.RdyMsg) {
			break
		}
		if timer.Get() >= timeout {
			return fmt.Errorf("while launching container, the ready messages does not appear after %v", timeout)
		}
	}
	return nil
}

func (r *Runner) Logs() (string, error) {
	cmdSlice := []string{"docker", "logs", r.ContainerName}
	outStr, errStr, err := exec.Run(cmdSlice)
	if err != nil {
		return "", fmt.Errorf(
			"Failed to inspect container '%s', stdout:%s stderr:%s error: %v",
			r.ContainerName,
			outStr,
			errStr,
			err,
		)
	}
	return outStr, nil
}

func (r *Runner) Inspect(filter string) (string, error) {
	cmdSlice := []string{"docker", "container", "inspect"}
	if len(filter) > 0 {
		cmdSlice = append(cmdSlice, "-f")
		cmdSlice = append(cmdSlice, filter)
	}
	cmdSlice = append(cmdSlice, r.ContainerName)
	outStr, errStr, err := exec.Run(cmdSlice)
	if err != nil {
		return "", fmt.Errorf(
			"Failed to inspect container '%s', stdout:%s stderr:%s error: %v",
			r.ContainerName,
			outStr,
			errStr,
			err,
		)
	}
	return outStr, nil
}

func (r *Runner) GetGatewayIP() (string, error) {
	ip, err := r.Inspect("{{range.NetworkSettings.Networks}}{{.Gateway}}{{end}}")
	if err != nil {
		return "", fmt.Errorf("failed to get container '%s' gateway IP, error: %v", r.ContainerName, err)
	}
	return ip, nil
}

func (r *Runner) Stop() error {
	cmdSlice := []string{"docker", "stop", r.ContainerName}
	outStr, errStr, err := exec.Run(cmdSlice)
	if err != nil {
		return fmt.Errorf("Could not run '%v', stdout=%s stderr=%s, error: %v", cmdSlice, outStr, errStr, err)
	}
	return nil
}

func (r *Runner) Wait() error {
	return r.Cmd.Wait()
}

// getPort returns the port of the string in the form of [IPv4|IPv6]:<port>
func getPort(addrPortStr string) (int, error) {
	components := strings.Split(addrPortStr, ":")
	if len(components) < 2 {
		return -1, fmt.Errorf("address:port line '%s' has less than 2 components separated by ':'", addrPortStr)
	}
	portStr := components[len(components)-1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return -1, fmt.Errorf("Port string '%s' is not a valid number", portStr)
	}
	return port, nil
}

// GetExposedPort return the exposed host port for the input containerPort.
func (r *Runner) GetExposedPort(containerPort int) (int, error) {
	dockerPortCmdSlice := []string{"docker", "port", r.ContainerName, strconv.Itoa(containerPort)}
	outStr, outErr, err := exec.Run(dockerPortCmdSlice)
	if err != nil {
		return -1, fmt.Errorf(
			"Failed to run docker command to get exposed port, stdout: %s, stderr: %s, error: %v",
			outStr,
			outErr,
			err,
		)
	}
	lines := strings.Split(outStr, "\n")
	if len(lines) == 0 {
		return -1, fmt.Errorf("Not output from docker Command '%v'", dockerPortCmdSlice)
	}
	for _, line := range lines {
		port, err := getPort(line)
		if err == nil {
			return port, err
		}
	}
	return -1, fmt.Errorf("Could not get exposed port for %d, output: %s", containerPort, outStr)
}
