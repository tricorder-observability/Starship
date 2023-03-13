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
	"net"

	"github.com/gin-gonic/gin"

	"github.com/tricorder/src/utils/log"

	swagfiles "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"

	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/pg"
)

type Config struct {
	Listen          net.Listener
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
	Standalone      bool
}

func StartHTTPService(cfg Config, pgClient *pg.Client) {
	log.Infof("Starting API Server's http service ...")

	grafana.InitGrafanaConfig(cfg.GrafanaURL, cfg.GrafanaUserName, cfg.GrafanaUserPass)

	grafanaManager := NewGrafanaManagement()
	err := grafanaManager.InitGrafanaAPIToken()
	if err != nil {
		msg := fmt.Sprintf("Failed to initialize Grafana API token, error: %v", err)
		if cfg.Standalone {
			log.Warnf(msg)
		} else {
			log.Fatalf(msg)
		}
	}

	mgr := ModuleManager{
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
	apiRoot.POST(api.CREATE_MODULE, mgr.createModuleHttp)
	apiRoot.GET(api.DELETE_MODULE, mgr.deleteModuleHttp)
	apiRoot.GET(api.LIST_AGENT, mgr.listAgentHttp)
	apiRoot.GET(api.LIST_MODULE, mgr.listModuleHttp)
	apiRoot.POST(api.DEPLOY_MODULE, mgr.deployModuleHttp)
	apiRoot.POST(api.UNDEPLOY_MODULE, mgr.undeployModuleHttp)

	router.GET("/swagger/*any", ginswag.WrapHandler(swagfiles.Handler))

	log.Infof("Listening on %s ...", cfg.Listen.Addr().String())
	_ = router.RunListener(cfg.Listen)
}
