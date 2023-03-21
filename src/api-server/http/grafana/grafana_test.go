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

package grafana

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	grafanaTest "github.com/tricorder/src/testing/grafana"
)

// Tests that auth token can be created on Grafana
func TestAuthToken(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cleanerFn, grafanaURL, err := grafanaTest.LaunchContainer()
	require.Nil(err)
	defer func() {
		// Have to write in this form
		// defer assert.Nil(cleanerFn())
		// causes cleanerFn() be invoked immediately, not really deferred.
		assert.Nil(cleanerFn())
	}()

	log.Println("grafana url:" + grafanaURL)
	config := NewConfig(grafanaURL, "admin", "admin")
	authToken := NewAuthToken(config)
	require.NotNil(authToken)

	token, err := authToken.GetToken("/api/dashboards/db")
	assert.Nil(err)

	dashboard := NewDashboard(config)
	assert.NotNil(dashboard)

	result, err := dashboard.CreateDashboard(token.Key, "APIServer1", "uid")
	assert.Nil(err)

	assert.Equal("success", result.Status)

	json, err := dashboard.GetDetailAsJSON(result.UID)
	assert.Nil(err)
	assert.Contains(json, `"title":"APIServer1"`)

	datasourceToken, err := authToken.GetToken("/api/datasources")
	require.Nil(err)
	assert.Nil(createDatasource(config, datasourceToken.Key))
}

func createDatasource(config Config, token string) error {
	ds := NewDatasource(config)
	const name = "MySQLTEST"

	if ds == nil {
		return fmt.Errorf("failed to create datasource")
	}

	_, err := ds.CreateDatasource(token, name,
		"localhost:5432", "postgres", "123456", "test")
	if err != nil {
		return fmt.Errorf("create grafana datasource error:%v", err)
	}
	return nil
}

// Tests that we can initialize API token on Grafana.
func TestInitGrafanaAPIToken(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cleanerFn, grafanaURL, err := grafanaTest.LaunchContainer()
	require.Nil(err)
	defer func() {
		assert.Nil(cleanerFn())
	}()

	config := NewConfig(grafanaURL, "admin", "admin")
	grafanaManager := NewGrafanaManagement(config)
	assert.Nil(grafanaManager.InitGrafanaAPIToken())
}
