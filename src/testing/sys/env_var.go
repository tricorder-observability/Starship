package sys

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetEnvVars() map[string]string {
	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		varName := pair[0]
		varValue := ""
		if len(pair) > 1 {
			varValue = pair[1]
		}
		log.Debugf("%s=%s", varName, varValue)
		envVars[varName] = varValue
	}
	return envVars
}
