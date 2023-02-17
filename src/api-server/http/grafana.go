package http

import (
	"fmt"
	"time"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"

	// Load sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type GrafanaManagement struct {
	GrafanaAPIKey dao.GrafanaAPIKey
}

var (
	dashboardAPIURL      = "/api/dashboards/db"
	datasourceAPIURL     = "/api/datasources"
	dashboardAPIURLName  = "dashboardAPIURL"
	datasourceAPIURLName = "datasourceAPIURL"
)

// TODO(yzhao): Rename apiKey to apiPath
func (g *GrafanaManagement) getGrafanaKey(apiKey, apiName string) (*dao.GrafanaAPIKeyGORM, error) {
	grafanaAPIKey, _ := g.GrafanaAPIKey.QueryByAPIKey(apiKey)
	if grafanaAPIKey == nil {
		authToken := grafana.NewAuthToken()
		token, err := authToken.GetToken(apiKey)
		if err != nil {
			return nil, fmt.Errorf("get grafana api token error %v", err)
		}
		if len(token.Key) > 0 {
			grafanaAPIKey = &dao.GrafanaAPIKeyGORM{
				Name:       apiName,
				APIKEY:     apiKey,
				AuthValue:  token.Key,
				CreateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			err = g.GrafanaAPIKey.SaveGrafanaAPI(grafanaAPIKey)
			if err != nil {
				return nil, fmt.Errorf("create grafana api token error %v", err)
			}
		}
	}
	return grafanaAPIKey, nil
}

func (g *GrafanaManagement) InitGrafanaAPIToken() error {
	_, err := g.getGrafanaKey(datasourceAPIURL, datasourceAPIURLName)
	if err != nil {
		return nil
	}
	_, err = g.getGrafanaKey(dashboardAPIURL, dashboardAPIURLName)
	if err != nil {
		return nil
	}
	time.Sleep(3 * time.Second)
	return nil
}
