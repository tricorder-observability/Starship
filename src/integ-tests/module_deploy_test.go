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

package integ_tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/sys"
)

// Tests that the http service can handle request
func TestGetDeployReqForModule(t *testing.T) {
	assert := assert.New(t)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	addrStr := fmt.Sprintf(":%d", *mgmtUIPort)
	listener, err := net.Listen(sys.TCP, addrStr)
	if err != nil {
		return errors.Wrap("starting http server", "listen", err)
	}
	config := http.Config{
		Listen:          listener,
		GrafanaURL:      *moduleGrafanaURL,
		GrafanaUserName: *moduleGrafanaUserName,
		GrafanaUserPass: *moduleGrafanaUserPassword,
		DatasourceName:  *moduleDatasourceName,
		DatasourceUID:   *moduleDatasourceUID,
		Module:          moduleDao,
		NodeAgent:       nodeAgentDao,
		ModuleInstance:  moduleInstanceDao,
		WaitCond:        waitCond,
		GLock:           gLock,
		Standalone:      *standalone,
	}
	return http.StartHTTPService(config, pgClient)
}
