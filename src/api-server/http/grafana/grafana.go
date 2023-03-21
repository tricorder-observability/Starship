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

package grafana

import (
	"github.com/tricorder/src/utils/errors"

	// Load sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type GrafanaManagement struct {
	config Config
}

const DashboardAPIURL = "/api/dashboards/db"

var grafanaAPIKeyMap = make(map[string]string)

func NewGrafanaManagement(config Config) GrafanaManagement {
	return GrafanaManagement{
		config: config,
	}
}

func (g *GrafanaManagement) GetGrafanaKey(apiPath string) (string, error) {
	if token, isExist := grafanaAPIKeyMap[apiPath]; isExist {
		return token, nil
	}

	authToken := NewAuthToken(g.config)
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
	grafanaAPIKeyMap[apiPath] = grafanaToken.Key

	return grafanaToken.Key, nil
}

func (g *GrafanaManagement) InitGrafanaAPIToken() error {
	_, err := g.GetGrafanaKey(DashboardAPIURL)
	if err != nil {
		return err
	}
	return nil
}
