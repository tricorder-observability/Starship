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

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	pb "github.com/tricorder/src/api-server/pb"
	testutils "github.com/tricorder/src/testing/bazel"
	grafanatest "github.com/tricorder/src/testing/grafana"
	pgclienttest "github.com/tricorder/src/testing/pg"
)

var cm = ModuleManager{DatasourceUID: "test"}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)
	cm.Module = dao.ModuleDao{
		Client: sqliteClient,
	}
	cm.GrafanaClient = NewGrafanaManagement()
	return router
}

// test upload wasm file and create wasm uid
func TestModuleManager(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	r := SetUpRouter()

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	grafana.InitGrafanaConfig(grafanaURL, "admin", "admin")

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	cm.PGClient = pgClient

	r.GET("/api/listModule", cm.listModuleHttp)
	req, _ := http.NewRequest("GET", "/api/listModule?fields=id,name,desire_state", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resultStr := w.Body.String()
	fmt.Printf("list module: %s", resultStr)
	assert.Equal(true, strings.Contains(resultStr, "Success"))

	wasmUid := "test_wasm_uid"
	modulID := AddModule(t, wasmUid, r)

	deployModule(t, modulID, r)

	unDeployModule(t, modulID, r)

	deleteModule(t, modulID, r)
}

func AddModule(t *testing.T, wasmUid string, r *gin.Engine) string {
	moduleName := "test_module"
	moduleBody := fmt.Sprintf(`{
		"name": "%s",
		"wasm":{
			"code": "",
			"fn_name":"copy_input_to_output",
			"output_schema":{
				"name":"test_tabel_name",
				"fields":[
					{
						"name":"data",
						"type": 5
					}
				]
			}
		},
		"ebpf":{
			"code": "",
			"perf_buffer_name":"events",
			"probes":[
				{
					"target":"",
					"entry":"sample_json",
					"return":""
				}
			]
		}
	}`, moduleName)
	jsonData := []byte(moduleBody)
	r.POST("/api/addModule", cm.createModuleHttp)
	req, _ := http.NewRequest("POST", "/api/addModule", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resultStr := w.Body.String()
	fmt.Printf("add module: %s", resultStr)
	// check http result
	assert.Equal(t, true, strings.Contains(resultStr, "success"))

	// check db result
	moduleResult, err := cm.Module.QueryByName(moduleName)
	if err != nil {
		t.Errorf("query module by name error:%v", err)
	}
	// check whether the name in the database is moduleName
	assert.Equal(t, true, moduleName == moduleResult.Name)
	assert.Equal(t, true, int(pb.DeploymentState_CREATED) == moduleResult.DesireState)
	return moduleResult.ID
}

// Tests that createModuleHttp failed if the input data fields are empty.
func TestCreateModuleEmptyDataFields(t *testing.T) {
	assert := assert.New(t)

	moduleName := "test_module"
	moduleBody := fmt.Sprintf(`{
		"name": "%s",
		"wasm":{
			"code": "",
			"fn_name":"copy_input_to_output",
			"output_schema":{
				"name":"test_tabel_name",
				"fields":[]
			}
		},
		"ebpf":{
			"code": "",
			"perf_buffer_name":"events",
			"probes":[
				{
					"target":"",
					"entry":"sample_json",
					"return":""
				}
			]
		}
	}`, moduleName)

	jsonData := []byte(moduleBody)
	r := SetUpRouter()
	r.POST("/api/createModule", cm.createModuleHttp)
	req, _ := http.NewRequest("POST", "/api/createModule", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(`{"code":500,"message":"input data fields cannot be empty"}`, w.Body.String())
}

func deleteModule(t *testing.T, modulID string, r *gin.Engine) {
	r.GET("/api/deleteModule", cm.deleteModuleHttp)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/deleteModule?id=%s", modulID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resultStr := w.Body.String()
	fmt.Printf("delete module: %s", resultStr)
	assert.Equal(t, true, strings.Contains(resultStr, "Success"))

	resultModule, _ := cm.Module.QueryByID(modulID)
	if resultModule != nil {
		t.Errorf("delete module by id error:%v", resultModule)
	}
}

func unDeployModule(t *testing.T, modulID string, r *gin.Engine) {
	r.GET("/api/undeployModule", cm.undeployModuleHttp)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/undeployModule?id=%s", modulID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resultStr := w.Body.String()
	fmt.Printf("un deploy module: %s", resultStr)
	assert.Equal(t, true, strings.Contains(resultStr, "un-deploy success"))

	// check code's status
	resultModule, err := cm.Module.QueryByID(modulID)
	if err != nil {
		t.Errorf("query module by id error:%v", err)
	}
	assert.Equal(t, int(pb.DeploymentState_TO_BE_UNDEPLOYED), resultModule.DesireState)
}

func deployModule(t *testing.T, modulID string, r *gin.Engine) {
	r.GET("/api/deployModule", cm.deployModuleHttp)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/deployModule?id=%s", modulID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resultStr := w.Body.String()

	assert.Equal(t, true, strings.Contains(resultStr, "prepare to deploy"))

	var deployResult DeployModuleResp
	err := json.Unmarshal([]byte(resultStr), &deployResult)
	if err != nil {
		t.Errorf("deploy module error:%v", err)
	}

	// check module's status
	moduleResult, err := cm.Module.QueryByID(modulID)
	if err != nil {
		t.Errorf("query module by id error:%v", err)
	}
	assert.Equal(t, int(pb.DeploymentState_TO_BE_DEPLOYED), moduleResult.DesireState)

	// check grafana dashboard create result
	ds := grafana.NewDashboard()
	detailResult, err := ds.GetDashboardDetail(deployResult.UID)
	if err != nil {
		t.Errorf("get grafana dashboard detail error:%v", err)
	}
	temp := fmt.Sprint(detailResult)
	assert.Equal(t, true, strings.Contains(temp, deployResult.UID))

	// check create postgres schema result
	const moduleDataTableNamePrefix = "tricorder_module_"
	err = cm.PGClient.CheckTableExist(moduleDataTableNamePrefix + moduleResult.ID)
	if err != nil {
		t.Errorf("check postgress table exist error:%v", err)
	}
}
