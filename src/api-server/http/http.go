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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	http_utils "github.com/tricorder/src/utils/http"
	"github.com/tricorder/src/utils/pg"
)

type Config struct {
	Port            int
	GrafanaURL      string
	GrafanaUserName string
	GrafanaUserPass string
	DatasourceName  string
	DatasourceUID   string
	Module          dao.Module
	GrafanaAPIKey   dao.GrafanaAPIKey
}

func StartHTTPService(cfg Config, pgClient *pg.Client) {
	log.Infof("Starting API Server's http service ...")

	grafana.InitGrafanaConfig(cfg.GrafanaURL, cfg.GrafanaUserName, cfg.GrafanaUserPass)

	grafanaManager := GrafanaManagement{
		GrafanaAPIKey: cfg.GrafanaAPIKey,
	}
	err := grafanaManager.InitGrafanaAPIToken()
	if err != nil {
		log.Fatalf("Failed to initialize Grafana API token, error: %v", err)
	}

	cm := ModuleManager{
		DatasourceUID: cfg.DatasourceUID,
		GrafanaClient: grafanaManager,
		Module:        cfg.Module,
		PGClient:      pgClient,
	}
	router := gin.Default()

	router.Use(Cors()).Use(GlobalExceptionWare)

	// TODO: Use swagger to define these APIs.
	api := router.Group(fmt.Sprintf("/%s", http_utils.API_ROOT))
	{
		api.POST(fmt.Sprintf("/%s", http_utils.ADD_CODE), cm.createModule)
		api.GET(fmt.Sprintf("/%s", http_utils.DELETE_MODULE), cm.deleteCode)
		api.GET(fmt.Sprintf("/%s", http_utils.LIST_CODE), cm.listCode)
		api.POST(fmt.Sprintf("/%s", http_utils.DEPLOY), cm.deployCode)
		api.POST(fmt.Sprintf("/%s", http_utils.UN_DEPLOY), cm.undeployCode)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Infof("Listening on %s ...", addr)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = router.Run(addr)
}
