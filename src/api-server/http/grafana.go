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

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"

	// Load sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type GrafanaManagement struct {
	GrafanaAPIKey    dao.GrafanaAPIKey
	GrafanaAPIKeyMap map[string]string
}

var (
	dashboardAPIURL     = "/api/dashboards/db"
	dashboardAPIURLName = "dashboardAPIURL"
)

// TODO(yzhao): Rename apiKey to apiPath
func (g *GrafanaManagement) getGrafanaKey(apiPath, apiName string) (string, error) {
	if token, isExist := g.GrafanaAPIKeyMap[apiPath]; isExist {
		return token, nil
	}

	authToken := grafana.NewAuthToken()
	allGrafanaAPIKey, err := authToken.GetAllGrafanaAPIKey()
	if err != nil {
		return "", fmt.Errorf("get grafana all api token list error %v", err)
	}
	for _, value := range allGrafanaAPIKey {
		err = authToken.RemoveGrafanaAPIKeyById(value.ID)
		if err != nil {
			return "", fmt.Errorf("remove grafana api token error %v", err)
		}
	}

	grafanaToken, err := authToken.GetToken(apiPath)
	if err != nil {
		return "", fmt.Errorf("get grafana api token error %v", err)
	}
	if len(g.GrafanaAPIKeyMap) == 0 {
		g.GrafanaAPIKeyMap = make(map[string]string)
	}
	g.GrafanaAPIKeyMap[apiPath] = grafanaToken.Key

	return grafanaToken.Key, nil
}

func (g *GrafanaManagement) InitGrafanaAPIToken() error {
	_, err := g.getGrafanaKey(dashboardAPIURL, dashboardAPIURLName)
	if err != nil {
		return nil
	}
	return nil
}
