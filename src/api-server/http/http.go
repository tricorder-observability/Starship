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
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/tricorder/src/utils/log"

	swagfiles "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/api-server/http/grafana"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/pg"
)

type Config struct {
	Port            int
	GrafanaURL      string
	GrafanaUserName string
	GrafanaUserPass string
	DatasourceName  string
	DatasourceUID   string
	Module          dao.ModuleDao
	NodeAgent       dao.NodeAgentDao
	ModuleInstance  dao.ModuleInstanceDao
	GLock           *lock.Lock
	WaitCond        *cond.Cond
}

func StartHTTPService(cfg Config, pgClient *pg.Client) {
	log.Infof("Starting API Server's http service ...")

	grafana.InitGrafanaConfig(cfg.GrafanaURL, cfg.GrafanaUserName, cfg.GrafanaUserPass)

	grafanaManager := NewGrafanaManagement()
	err := grafanaManager.InitGrafanaAPIToken()
	if err != nil {
		log.Fatalf("Failed to initialize Grafana API token, error: %v", err)
	}

	cm := ModuleManager{
		DatasourceUID:  cfg.DatasourceUID,
		GrafanaClient:  grafanaManager,
		Module:         cfg.Module,
		NodeAgent:      cfg.NodeAgent,
		ModuleInstance: cfg.ModuleInstance,
		PGClient:       pgClient,
		gLock:          cfg.GLock,
		waitCond:       cfg.WaitCond,
	}
	router := gin.Default()

	router.Use(Cors()).Use(GlobalExceptionWare)

	apiRoot := router.Group(api.ROOT)
	apiRoot.POST(api.CREATE_MODULE, cm.createModuleHttp)
	apiRoot.GET(api.DELETE_MODULE, cm.deleteModuleHttp)
	apiRoot.GET(api.LIST_MODULE, cm.listModuleHttp)
	apiRoot.POST(api.DEPLOY_MODULE, cm.deployModuleHttp)
	apiRoot.POST(api.UNDEPLOY_MODULE, cm.undeployModuleHttp)

	router.GET("/swagger/*any", ginswag.WrapHandler(swagfiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Infof("Listening on %s ...", addr)
	_ = router.Run(addr)
}
