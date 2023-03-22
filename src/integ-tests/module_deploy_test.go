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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/wasm"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/sys"

	"github.com/tricorder/src/testing/bazel"
	grafanatest "github.com/tricorder/src/testing/grafana"
	pgclienttest "github.com/tricorder/src/testing/pg"
)

// Tests that the http service can handle request
func TestGetDeployReqForModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	listener, httpSrvAddr, err := sys.ListenTCP(0)
	require.NoError(err)
	assert.Regexp(`\[::\]:[0-9]+`, httpSrvAddr.String())

	dirPath := bazel.CreateTmpDir()

	sqliteClient, err := dao.InitSqlite(dirPath)
	require.NoError(err)

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}
	nodeAgentDao := dao.NodeAgentDao{
		Client: sqliteClient,
	}
	moduleInstanceDao := dao.ModuleInstanceDao{
		Client: sqliteClient,
	}
	config := http.Config{
		Listen:          listener,
		GrafanaURL:      grafanaURL,
		GrafanaUserName: "admin",
		GrafanaUserPass: "admin",
		DatasourceName:  "TimescaleDB-Tricorder",
		DatasourceUID:   "timescaledb_tricorder",
		Module:          moduleDao,
		NodeAgent:       nodeAgentDao,
		ModuleInstance:  moduleInstanceDao,
		WaitCond:        cond.NewCond(),
		GLock:           lock.NewLock(),
	}

	wasiCompiler := wasm.NewWASICompilerWithDefaults()
	go func() {
		err := http.StartHTTPService(config, pgClient, wasiCompiler)
		require.NoError(err)
	}()
}
