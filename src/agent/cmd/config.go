package main

import (
	"fmt"
	"strings"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/sys"
)

type config struct {
	// The name of the Kubernetes node on which this agent is running.
	// This is used to distinguish the currently-running agent and the terminated one after agent restarts.
	nodeName string

	// The ID of the Pod in which this agent is running.
	// This is used by API Server to check if a agent is already terminated, by observing whether or not this pod has been
	// deleted.
	podID string
}

const (
	NODE_NAME string = "NODE_NAME"
	POD_ID    string = "POD_ID"
)

var requiredEnvVarNames []string = []string{NODE_NAME, POD_ID}

func checkRequiredEnvVarsAreDefined() error {
	var missingVarNames []string
	envVars := sys.EnvVars()
	for _, n := range requiredEnvVarNames {
		val, found := envVars[n]
		if !found || len(val) != 0 {
			missingVarNames = append(missingVarNames, n)
		}
	}
	if len(missingVarNames) > 0 {
		return fmt.Errorf("required env vars [%s], missing [%s]", strings.Join(requiredEnvVarNames, ", "),
			strings.Join(missingVarNames, ", "))
	}
	return nil
}

func newConfig() (*config, error) {
	err := checkRequiredEnvVarsAreDefined()
	if err != nil {
		return *errors.Wrap("newing agent config", "check env vars", err)
	}
	envVars := sys.EnvVars()

}
