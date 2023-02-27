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
	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	"github.com/tricorder/src/utils/errors"

	// Load sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type GrafanaManagement struct {
	grafanaAPIKey    dao.GrafanaAPIKey
	grafanaAPIKeyMap map[string]string
}

var (
	dashboardAPIURL = "/api/dashboards/db"
)

func NewGrafanaManagement(grafananAPIKey dao.GrafanaAPIKey) GrafanaManagement {
	grafanaManager := GrafanaManagement{
		grafanaAPIKey: grafananAPIKey,
	}
	grafanaManager.grafanaAPIKeyMap = make(map[string]string)
	return grafanaManager
}

func (g *GrafanaManagement) getGrafanaKey(apiPath string) (string, error) {
	if token, isExist := g.grafanaAPIKeyMap[apiPath]; isExist {
		return token, nil
	}

	authToken := grafana.NewAuthToken()
	allGrafanaAPIKey, err := authToken.GetAllGrafanaAPIKey()
	if err != nil {
		return "", errors.Wrap("get grafana all api token list", "load", err)
	}
	for _, value := range allGrafanaAPIKey {
		err = authToken.RemoveGrafanaAPIKeyById(value.ID)
		if err != nil {
			return "", errors.Wrap("remove grafana api token", "load", err)
		}
	}

	grafanaToken, err := authToken.GetToken(apiPath)
	if err != nil {
		return "", errors.Wrap("get grafana api token", "load", err)
	}
	g.grafanaAPIKeyMap[apiPath] = grafanaToken.Key

	return grafanaToken.Key, nil
}

func (g *GrafanaManagement) InitGrafanaAPIToken() error {
	_, err := g.getGrafanaKey(dashboardAPIURL)
	if err != nil {
		return err
	}
	return nil
}
