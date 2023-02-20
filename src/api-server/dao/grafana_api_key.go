// Copyright (C) 2023  tricorder-observability
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

package dao

import (
	"fmt"

	"github.com/tricorder/src/utils/sqlite"
)

// GrafanaAPIKeyGORM sqlite gorm storage object
type GrafanaAPIKeyGORM struct {
	ID         int    `gorm:"'id' primarykey"`
	Name       string `gorm:"name"`
	CreateTime string `gorm:"create_time"`
	// APIKEY is the api path, used for querying with the sqlite db.
	APIKEY    string `gorm:"column:api_key"`
	AuthValue string `gorm:"auth_value"`
}

func (GrafanaAPIKeyGORM) TableName() string {
	// TODO: could be renamed to grafana_api_key
	return "grafana_api"
}

// TODO(zhihui): Rename to GrafanaAPIKeyDao
type GrafanaAPIKey struct {
	Client *sqlite.ORM
}

func (g *GrafanaAPIKey) SaveGrafanaAPI(grafanaAPI *GrafanaAPIKeyGORM) error {
	result := g.Client.Engine.Create(grafanaAPI)
	return result.Error
}

func (g *GrafanaAPIKey) QueryByAPIKey(apiKEY string) (*GrafanaAPIKeyGORM, error) {
	grafanaAPI := &GrafanaAPIKeyGORM{}
	result := g.Client.Engine.Where(&GrafanaAPIKeyGORM{APIKEY: apiKEY}).First(grafanaAPI)
	if result.Error != nil {
		return nil, fmt.Errorf("query grafana api by api key error:%v", result.Error)
	}
	return grafanaAPI, nil
}

func (g *GrafanaAPIKey) QueryByID(id int) (*GrafanaAPIKeyGORM, error) {
	grafanaAPI := &GrafanaAPIKeyGORM{}
	result := g.Client.Engine.Where(&GrafanaAPIKeyGORM{ID: id}).First(grafanaAPI)
	if result.Error != nil {
		return nil, fmt.Errorf("query grafana api by id error:%v", result.Error)
	}
	return grafanaAPI, nil
}
