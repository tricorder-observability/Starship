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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	pb "github.com/tricorder/src/api-server/pb"
	modulepb "github.com/tricorder/src/pb/module"
	commonpb "github.com/tricorder/src/pb/module/common"
	"github.com/tricorder/src/utils/channel"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/uuid"
)

// ModuleManager provides APIs to manage eBPF+WASM module received from the management Web UI.
type ModuleManager struct {
	DatasourceUID string
	Module        dao.Module
	GrafanaClient GrafanaManagement
	PGClient      *pg.Client
}

// ShowAccount godoc
// @Summary      Add module
// @Description  Create Module
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param			   module	body	module.Module	true	"Add module"
// @Success      200  {object}  module.Module
// @Router       /api/addCode [post]
func (mgr *ModuleManager) createModule(c *gin.Context) {
	var body modulepb.Module

	err := c.ShouldBind(&body)
	if err != nil {
		log.Errorf("while creating module, failed to bind request body to module structure, error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "Request Error: " + err.Error()})
		return
	}

	m, _ := mgr.Module.QueryByName(body.Name)
	if m != nil && len(m.Name) > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": fmt.Sprintf("Name '%s' already exists", body.Name)})
		return
	}

	ebpfProbes, err := json.Marshal(body.Ebpf.Probes)
	if err != nil {
		log.Errorf("while creating module, failed to marshal ebpf probespecs, error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "request error: " + err.Error()})
		return
	}

	schemaAttr, err := json.Marshal(body.Wasm.OutputSchema.Fields)
	if err != nil {
		log.Errorf("while creating module, failed to marshal wasm output fields, error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "request error: " + err.Error()})
		return
	}

	mod := &dao.ModuleGORM{
		ID:                 strings.Replace(uuid.New(), "-", "_", -1),
		Name:               body.Name,
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		Status:             int(pb.DeploymentState_CREATED),
		Ebpf:               body.Ebpf.Code,
		EbpfFmt:            int(body.Ebpf.Fmt),
		EbpfLang:           int(body.Ebpf.Lang),
		EbpfPerfBufferName: body.Ebpf.PerfBufferName,
		EbpfProbes:         string(ebpfProbes),
		Wasm:               body.Wasm.Code,
		SchemaAttr:         string(schemaAttr),
		Fn:                 body.Wasm.FnName,
		WasmFmt:            int(body.Wasm.Fmt),
		WasmLang:           int(body.Wasm.Lang),
	}

	mod.SchemaName = fmt.Sprintf("%s_%s", "tricorder_code", mod.ID)

	err = mgr.Module.SaveCode(mod)
	if err != nil {
		log.Errorf("save code module error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "Save code module error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "create success, module id: " + mod.ID})
}

// ShowAccount godoc
// @Summary      List all moudle
// @Description  Create Module
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param			   fields	 query	string	false  "query field search like 'id,name,createTime'"
// @Success      200  {array}  module.Module
// @Router       /api/listCode [get]
func (mgr *ModuleManager) listCode(c *gin.Context) {
	var resultList []dao.ModuleGORM
	var err error
	// ?fields=id,name,status
	fields, exist := c.GetQuery("fields")
	if exist {
		resultList, err = mgr.Module.ListCode(fields)
	} else {
		resultList, err = mgr.Module.ListCode()
	}

	if err != nil {
		log.Errorf("Failed to list code, error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "Query Error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Success", "data": resultList})
}

// ShowAccount godoc
// @Summary      Delete module by id
// @Description  Create Module
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param			   id	  query		  string	true	"delete module id"
// @Success      200  {object}  module.Module
// @Router       /api/deleteCode [get]
func (mgr *ModuleManager) deleteCode(c *gin.Context) {
	id, exist := c.GetQuery("id")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "id does not exist"})
		return
	}
	err := mgr.Module.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "delete error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Success"})
}

// ShowAccount godoc
// @Summary      Deploy module
// @Description  Create Module
// @Tags         module
// @Accept       json
// @Produce      json
// @Param			   id	  query		  string	true	"deploy module id"
// @Success      200  {object}  module.Module
// @Router       /api/deployCode [post]
func (mgr *ModuleManager) deployCode(c *gin.Context) {
	id, exist := c.GetQuery("id")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "id cannot be empty"})
		return
	}
	// Check whether the code exists
	code, err := mgr.Module.QueryByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "module does not exist"})
		return
	}

	err = mgr.createPGTable(code)
	if err != nil {
		log.Error("Failed to create PG table")
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "create schema error: " + err.Error()})
		return
	}
	log.Info("Created postgres table")

	uid, err := mgr.createGrafanaDashboard(code.ID)
	if err != nil {
		log.Error("Failed to create Grafana dashboard")
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "create dashboard error"})
		return
	}
	log.Infof("Created Grafana dashboard with UID: %s", uid)

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "prepare to deploy module, id: " + id})
}

// ShowAccount godoc
// @Summary      Undeploy module
// @Description  Create Module
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param			   id	  query		 string	 true	 "undeploy module id"
// @Success      200  {object}  module.Module
// @Router       /api/undeployCode [post]
func (mgr *ModuleManager) undeployCode(c *gin.Context) {
	id, exist := c.GetQuery("id")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "id cannot be empty"})
		return
	}
	err := mgr.Module.UpdateStatusByID(id, int(pb.DeploymentState_TO_BE_UNDEPLOYED))
	if err != nil {
		log.Errorf("pre-undeploy code: [%s] failed: %s", id, err.Error())
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "undeploy error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "un-deploy success"})
}

// Generate schema name tricorder_code_{moduleID}
func getModuleDataTableName(id string) string {
	const moduleDataTableNamePrefix = "tricorder_code_"
	return moduleDataTableNamePrefix + id
}

// createPGTable creates a data table on the database that stores observability data.
// Agents can then write the data produced by the deployed eBPF+WASM module to this table.
func (mgr *ModuleManager) createPGTable(code *dao.ModuleGORM) error {
	var fields []*commonpb.DataField
	err := json.Unmarshal([]byte(code.SchemaAttr), &fields)
	if err != nil {
		return fmt.Errorf("while creating output data table for module '%s', "+
			"failed to unmarshal column schemas, error: %v", code.Name, err)
	}
	if len(fields) == 0 {
		log.Infof("Module '%s' has no data schema defined", code.Name)
		return nil
	}
	columns, err := DataFieldsToPGColumns(fields)
	if err != nil {
		return err
	}
	schema := pg.Schema{
		Name:    getModuleDataTableName(code.ID),
		Columns: columns,
	}
	err = mgr.PGClient.CreateTable(&schema)
	if err != nil {
		return fmt.Errorf("while creating output data table for module '%s', "+
			"failed to create the table, error: %v", schema.Name, err)
	}
	return nil
}

func (mgr *ModuleManager) createGrafanaDashboard(moduleID string) (string, error) {
	grafanaAPIKey, err := mgr.GrafanaClient.getGrafanaKey(dashboardAPIURL, dashboardAPIURLName)
	if err != nil {
		log.Println("deploy error, auth dashboary error", err)
		return "", err
	}

	ds := grafana.NewDashboard()
	result, err := ds.CreateDashboard(grafanaAPIKey.AuthValue, getModuleDataTableName(moduleID), mgr.DatasourceUID)
	if err != nil {
		log.Println("Create dashboard", err)
		return "", err
	}

	err = mgr.Module.UpdateStatusByID(moduleID, int(pb.DeploymentState_TO_BE_DEPLOYED))
	if err != nil {
		log.Errorf("pre-deploy code: [%s] failed: %s", moduleID, err.Error())
		return "", err
	}

	message := channel.DeployChannelModule{
		ID:     moduleID,
		Status: int(pb.DeploymentState_TO_BE_DEPLOYED),
	}
	channel.SendMessage(message)

	return result.UID, nil
}
