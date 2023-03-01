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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	pb "github.com/tricorder/src/api-server/pb"
	commonpb "github.com/tricorder/src/pb/module/common"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/uuid"
)

// ModuleManager provides APIs to manage eBPF+WASM module received from the management Web UI.
type ModuleManager struct {
	DatasourceUID  string
	Module         dao.ModuleDao
	NodeAgent      dao.NodeAgentDao
	ModuleInstance dao.ModuleInstanceDao
	GrafanaClient  GrafanaManagement
	gLock          *lock.Lock
	waitCond       *cond.Cond
	PGClient       *pg.Client
}

// createModuleHttp  godoc
// @Summary      Create module
// @Description  Store module data into SQLite database
// @Tags         module
// @Accept       json
// @Produce      json
// @Param			   module	body	CreateModuleReq	true	"Create module"
// @Success      200  {object}  CreateModuleResp
// @Router       /api/addModule [post]
func (mgr *ModuleManager) createModuleHttp(c *gin.Context) {
	var body CreateModuleReq
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": "Request Error: " + err.Error()})
		return
	}
	result := mgr.createModule(body)
	c.JSON(http.StatusOK, result)
}

func (mgr *ModuleManager) createModule(body CreateModuleReq) CreateModuleResp {
	var m *dao.ModuleGORM
	var err error
	err = mgr.gLock.ExecWithLock(func() error {
		m, _ = mgr.Module.QueryByName(body.Name)
		if m != nil && len(m.Name) > 0 {
			return fmt.Errorf("Name '%s' already exists", body.Name)
		}
		return nil
	})

	if err != nil {
		return CreateModuleResp{HTTPResp{
			Code:    500,
			Message: err.Error(),
		}}
	}

	ebpfProbes, err := json.Marshal(body.Ebpf.Probes)
	if err != nil {
		log.Errorf("while creating module, failed to marshal ebpf probespecs, error: %v", err)
		return CreateModuleResp{HTTPResp{
			Code:    500,
			Message: "request error: " + err.Error(),
		}}
	}

	if len(body.Wasm.OutputSchema.Fields) == 0 {
		return CreateModuleResp{HTTPResp{
			Code:    500,
			Message: "input data fields cannot be empty",
		}}
	}

	schemaAttr, err := json.Marshal(body.Wasm.OutputSchema.Fields)
	if err != nil {
		return CreateModuleResp{HTTPResp{
			Code:    500,
			Message: "request error: " + err.Error(),
		}}
	}

	mod := &dao.ModuleGORM{
		ID:                 strings.Replace(uuid.New(), "-", "_", -1),
		Name:               body.Name,
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		DesireState:        int(pb.ModuleState_CREATED_),
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

	mod.SchemaName = fmt.Sprintf("%s_%s", "tricorder_module", mod.ID)

	err = mgr.gLock.ExecWithLock(func() error {
		err = mgr.Module.SaveModule(mod)
		if err != nil {
			log.Errorf("save module error: %v", err)
			return errors.New("Save module error: " + err.Error())
		}
		return nil
	})

	if err != nil {
		return CreateModuleResp{HTTPResp{
			Code:    500,
			Message: err.Error(),
		}}
	}

	return CreateModuleResp{HTTPResp{
		Code:    200,
		Message: "create success, module id: " + mod.ID,
	}}
}

// listModuleHttp godoc
// @Summary      List all moudle
// @Description  List all moudle
// @Tags         module
// @Accept       json
// @Produce      json
// @Param			   fields	 query	string	false  "query field search like 'id,name,createTime'"
// @Success      200  {object}  ListModuleResp
// @Router       /api/listModule [get]
func (mgr *ModuleManager) listModuleHttp(c *gin.Context) {
	fields, err := checkQuery(c, "fields")
	if err != nil {
		return
	}
	result := mgr.listModule(ListModuleReq{Fields: fields})
	c.JSON(http.StatusOK, result)
}

func (mgr *ModuleManager) listModule(req ListModuleReq) ListModuleResp {
	resultList, err := mgr.Module.ListModule(req.Fields)
	if err != nil {
		log.Errorf("Failed to list module, error: %v", err)
		return ListModuleResp{HTTPResp{
			Code:    500,
			Message: "Query Error: " + err.Error(),
		}, nil}
	}

	return ListModuleResp{HTTPResp{
		Code:    200,
		Message: "Success",
	}, resultList}
}

// deleteModuleHttp  godoc
// @Summary      Delete module
// @Description  Delete module by id
// @Tags         module
// @Accept       json
// @Produce      json
// @Param			   id	  query		  string	true	"delete module id"
// @Success      200  {object}   HTTPResp
// @Router       /api/deleteModule [get]
func (mgr *ModuleManager) deleteModuleHttp(c *gin.Context) {
	id, err := checkQuery(c, "id")
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, mgr.deleteModule(id))
}

func (mgr *ModuleManager) deleteModule(id string) DeleteModuleResp {
	err := mgr.Module.DeleteByID(id)
	if err != nil {
		return DeleteModuleResp{HTTPResp{
			Code:    500,
			Message: "delete error: " + err.Error(),
		}}
	}
	return DeleteModuleResp{HTTPResp{
		Code:    200,
		Message: "Success",
	}}
}

// deployModuleHttp godoc
// @Summary      Deploy module
// @Description  Deploy the specified module onto every agent in the cluster
// @Tags         module
// @Accept       json
// @Produce      json
// @Param			   id	  query		  string	true	"deploy module id"
// @Success      200  {object}  DeployModuleResp
// @Router       /api/deployModule [post]
func (mgr *ModuleManager) deployModuleHttp(c *gin.Context) {
	id, err := checkQuery(c, "id")
	if err != nil {
		return
	}
	result := mgr.deployModule(id)
	c.JSON(http.StatusOK, result)
}

func (mgr *ModuleManager) deployModule(id string) DeployModuleResp {
	var module *dao.ModuleGORM
	var err error
	// Check whether the module exists
	err = mgr.gLock.ExecWithLock(func() error {
		module, err = mgr.Module.QueryByID(id)
		if err != nil {
			return errors.New("query module error: " + err.Error())
		}
		return nil
	})

	if err != nil {
		return DeployModuleResp{
			HTTPResp{
				Code:    500,
				Message: err.Error(),
			},
			id,
		}
	}

	err = mgr.createPGTable(module)
	if err != nil {
		log.Error("Failed to create PG table")
		return DeployModuleResp{
			HTTPResp{
				Code:    500,
				Message: "create schema error: " + err.Error(),
			},
			"",
		}
	}
	log.Info("Created postgres table")

	uid, err := mgr.createGrafanaDashboard(module.ID)
	if err != nil {
		log.Error("Failed to create Grafana dashboard")

		return DeployModuleResp{
			HTTPResp{
				Code:    500,
				Message: "create dashboard error",
			},
			uid,
		}
	}

	log.Infof("Created Grafana dashboard with UID: %s", uid)

	err = mgr.gLock.ExecWithLock(func() error {
		err = mgr.Module.UpdateStatusByID(module.ID, int(pb.ModuleState_DEPLOYED))
		if err != nil {
			return errors.New("pre-deploy module: " + module.ID + "failed: " + err.Error())
		}
		return nil
	})

	if err != nil {
		return DeployModuleResp{
			HTTPResp{
				Code:    500,
				Message: err.Error(),
			},
			uid,
		}
	}

	mgr.waitCond.Broadcast()

	return DeployModuleResp{
		HTTPResp{
			Code:    200,
			Message: "prepare to deploy module, id: " + id,
		},
		uid,
	}
}

// undeployModuleHttp godoc
// @Summary      Undeploy module
// @Description  Undeploy the specified module from all agents in the cluster
// @Tags         module
// @Accept       json
// @Produce      json
// @Param			   id	  query		 string	 true	 "undeploy module id"
// @Success      200  {object}  HTTPResp
// @Router       /api/undeployModule [post]
func (mgr *ModuleManager) undeployModuleHttp(c *gin.Context) {
	id, err := checkQuery(c, "id")
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, mgr.undeployModule(id))
}

func (mgr *ModuleManager) undeployModule(id string) UndeployModuleResp {
	err := mgr.gLock.ExecWithLock(func() error {
		err := mgr.Module.UpdateStatusByID(id, int(pb.ModuleState_UNDEPLOYED))
		if err != nil {
			return errors.New("un-deploy module: " + id + "failed: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return UndeployModuleResp{HTTPResp{
			Code:    500,
			Message: err.Error(),
		}}
	}
	return UndeployModuleResp{HTTPResp{
		Code:    200,
		Message: "un-deploy success",
	}}
}

// Generate schema name tricorder_module_{moduleID}
func getModuleDataTableName(id string) string {
	const moduleDataTableNamePrefix = "tricorder_module_"
	return moduleDataTableNamePrefix + id
}

// createPGTable creates a data table on the database that stores observability data.
// Agents can then write the data produced by the deployed eBPF+WASM module to this table.
func (mgr *ModuleManager) createPGTable(module *dao.ModuleGORM) error {
	var fields []*commonpb.DataField
	err := json.Unmarshal([]byte(module.SchemaAttr), &fields)
	if err != nil {
		return fmt.Errorf("while creating output data table for module '%s', "+
			"failed to unmarshal column schemas, error: %v", module.Name, err)
	}
	if len(fields) == 0 {
		return fmt.Errorf("module data fields cannot be empty")
	}
	columns, err := DataFieldsToPGColumns(fields)
	if err != nil {
		return err
	}
	schema := pg.Schema{
		Name:    getModuleDataTableName(module.ID),
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
	grafanaAPIKey, err := mgr.GrafanaClient.getGrafanaKey(dashboardAPIURL)
	if err != nil {
		log.Println("deploy error, auth dashboary error", err)
		return "", err
	}

	ds := grafana.NewDashboard()
	result, err := ds.CreateDashboard(grafanaAPIKey, getModuleDataTableName(moduleID), mgr.DatasourceUID)
	if err != nil {
		log.Println("Create dashboard", err)
		return "", err
	}

	return result.UID, nil
}
