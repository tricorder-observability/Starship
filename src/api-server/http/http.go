package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	http_utils "github.com/tricorder/src/utils/http"
	"github.com/tricorder/src/utils/pg"

	_ "github.com/tricorder/src/utils/log"
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

	_ = router.Run(addr)
}
