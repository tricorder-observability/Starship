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
