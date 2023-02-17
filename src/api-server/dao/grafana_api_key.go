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
	APIKEY     string `gorm:"column:api_key"`
	AuthValue  string `gorm:"auth_value"`
}

func (GrafanaAPIKeyGORM) TableName() string {
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
