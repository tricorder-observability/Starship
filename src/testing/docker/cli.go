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

package testing

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/tricorder/src/utils/exec"
)

// CLI wraps the `docker` command
type CLI struct{}

// Load imports a local .tar file by calling `docker load -i`.
// Returns the URL of the imported image.
func (d *CLI) Load(tarFile string) string {
	cmdSlice := []string{"docker", "load", "-i", tarFile}
	outStr, outErr, err := exec.Run(cmdSlice)
	if err != nil {
		log.Fatalf(
			"Failed to run docker command to load tar file '%s', stdout: %s, stderr: %s, error: %v",
			tarFile,
			outStr,
			outErr,
			err,
		)
	}
	words := strings.Fields(outStr)
	return words[2]
}
