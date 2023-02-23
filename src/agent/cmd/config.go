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

package main

import (
	"fmt"
	"strings"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/sys"
)

// config records configuration entires for agent.
type config struct {
	// The name of the Kubernetes node on which this agent is running.
	// This is used to distinguish the currently-running agent and the terminated one after agent restarts.
	nodeName string

	// The ID of the Pod in which this agent is running.
	// This is used by API Server to check if a agent is already terminated, by observing whether or not this pod has been
	// deleted.
	podID string
}

// Environment variables are injected by [downwardAPI](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/)
// in the [helm-charts](https://github.com/tricorder-observability/helm-charts).
const (
	ENV_VAR_NODE_NAME string = "NODE_NAME"
	ENV_VAR_POD_ID    string = "POD_ID"
)

var requiredEnvVarNames = []string{ENV_VAR_NODE_NAME, ENV_VAR_POD_ID}

func checkRequiredEnvVarsAreDefined(envVars map[string]string) error {
	var missingVarNames []string
	for _, n := range requiredEnvVarNames {
		val, found := envVars[n]
		if !found || len(val) == 0 {
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
	envVars := sys.EnvVars()
	err := checkRequiredEnvVarsAreDefined(envVars)
	if err != nil {
		return nil, errors.Wrap("newing agent config", "check env vars", err)
	}
	c := new(config)
	c.nodeName = envVars[ENV_VAR_NODE_NAME]
	c.podID = envVars[ENV_VAR_POD_ID]
	return c, nil
}
