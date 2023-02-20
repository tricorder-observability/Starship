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

package dao

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
)

// test GrafanaAPIKey fun
// test SaveGrafanaAPI fun and check save result
// test QueryByID and QueryByAPIKey and check query result
func TestGrafanaAPIDAO(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	defer func() {
		assert.Nil(os.RemoveAll(dirPath))
	}()
	sqliteClient, _ := InitSqlite(dirPath)

	grafanaAPIDao := GrafanaAPIKey{
		Client: sqliteClient,
	}

	grafanaAPIKey := &GrafanaAPIKeyGORM{
		Name:       "testAPIName",
		APIKEY:     "APIKey",
		AuthValue:  "AuthTokenValue",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := grafanaAPIDao.SaveGrafanaAPI(grafanaAPIKey)
	if err != nil {
		t.Errorf("save grafana api err %v", err)
	}

	queryResult, err := grafanaAPIDao.QueryByID(grafanaAPIKey.ID)
	if err != nil {
		t.Errorf("query grafana api by id err %v", err)
	}
	if grafanaAPIKey.Name != queryResult.Name {
		t.Errorf("save grafana api by id error: not found name == %v", grafanaAPIKey.Name)
	}

	queryResult, err = grafanaAPIDao.QueryByAPIKey(grafanaAPIKey.APIKEY)
	if err != nil {
		t.Errorf("query grafana api by id err %v", err)
	}
	if grafanaAPIKey.APIKEY != queryResult.APIKEY {
		t.Errorf("save grafana api by api key error: not found apiKey == %v", grafanaAPIKey.APIKEY)
	}
}
