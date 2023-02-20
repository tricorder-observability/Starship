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

package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO(zhihui): Add more comments to explain the purpose of these types and APIs.

type Datasource struct {
	client http.Client
}

func NewDatasource() *Datasource {
	return &Datasource{
		client: http.Client{},
	}
}

func (g *Datasource) CreateDatasource(
	createDatabaseAuthKey, name, url, user, password, databases string,
) (*DatasourceResult, error) {
	bodyReq := DatasourceBody{
		Name:     name,
		Type:     "postgres",
		TypeName: "PostgreSQL",
		Access:   "proxy",
		URL:      url,
		Password: password,
		User:     user,
		Database: databases,
		SecureJSONData: SecureJSONData{
			Password: password,
		},
		JSONData: DataSourceJSONData{
			PostgresVersion: 1400,
			Sslmode:         "disable",
		},
	}

	bytesData, _ := json.Marshal(bodyReq)

	req, err := http.NewRequest("POST", CreateDatabaseURI, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+createDatabaseAuthKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var out DatasourceResult
		err = json.Unmarshal(body, &out)
		if err != nil {
			return nil, err
		}
		fmt.Println("CreateDatasource Result : ", out)
		return &out, nil
	}
	return nil, err
}

type DatasourceResult struct {
	ID         string         `json:"uid"`
	Message    string         `json:"message"`
	Name       string         `json:"name"`
	Datasource DatasourceBody `json:"datasource"`
}

type DatasourceBody struct {
	UID            string             `json:"uid"`
	Name           string             `json:"name"`
	Type           string             `json:"type"`
	TypeName       string             `json:"typeName"`
	Access         string             `json:"access"`
	URL            string             `json:"url"`
	Password       string             `json:"password"`
	User           string             `json:"user"`
	Database       string             `json:"database"`
	SecureJSONData SecureJSONData     `json:"secureJSONData"`
	JSONData       DataSourceJSONData `json:"jsonData"`
}

type SecureJSONData struct {
	Password string `json:"password"`
}

type DataSourceJSONData struct {
	PostgresVersion int    `json:"postgresVersion"`
	Sslmode         string `json:"sslmode"`
}
