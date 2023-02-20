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
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// TODO(zhihui): Add tests with a running Grafana instance in docker container
// https://github.com/tricorder-observability/starship/issues/534
// See this issue for some ideas

type Dashboard struct {
	client http.Client
}

func NewDashboard() *Dashboard {
	return &Dashboard{
		client: http.Client{},
	}
}

func (g *Dashboard) CreateDashboard(createDashBoardAuthKey, title, datasourceUID string) (*DashboardResult, error) {
	sql := "SELECT * FROM " + title + " LIMIT 50"
	targetList := [1]DashboardTargetData{{
		Format:       "table",
		MetricColumn: "none",
		RawQuery:     true,
		RawSQL:       sql,
		RefID:        "A",
		Table:        title,
		Hide:         false,
		Datasource: DashboardTargetDatasourceData{
			Type: "postgres",
			UID:  datasourceUID,
		},
	}}
	panelsObj := make([]DashboardPanelData, 1)
	panelsObj[0] = DashboardPanelData{
		Type:          "table",
		Title:         title,
		PluginVersion: "8.3.3",
		Targets:       targetList,
		GridPos: GridPos{
			X: 0,
			Y: 0,
			H: 15,
			W: 21,
		},
	}

	bodyReq := BodyData{
		Dashboard: DashboardData{
			Title:   title,
			Version: 1,
			Panels:  panelsObj,
			Refresh: "5s",
		},
	}

	bytesData, _ := json.Marshal(bodyReq)

	req, err := http.NewRequest("POST", CreateDashBoardURI, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+createDashBoardAuthKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var out DashboardResult
		err = json.Unmarshal(body, &out)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}
	return nil, err
}

func (g *Dashboard) AddDashboardPanel(
	createDashBoardAuthKey, uid, title, id string,
	panelsObj []DashboardPanelData,
) (*DashboardResult, error) {
	bodyReq := BodyData{
		Dashboard: DashboardData{
			ID:      id,
			UID:     uid,
			Title:   title,
			Version: 1,
			Panels:  panelsObj,
			Refresh: "5s",
		},
	}

	bytesData, _ := json.Marshal(bodyReq)

	req, err := http.NewRequest("POST", CreateDashBoardURI, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+createDashBoardAuthKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var out DashboardResult
		err = json.Unmarshal(body, &out)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}
	return nil, err
}

func (g *Dashboard) GetDashboardDetail(uid string) (*DashboardDetailResult, error) {
	req, err := http.NewRequest("GET", GetDashboardURI+uid, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(BasicAuth)))

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var out DashboardDetailResult
		err = json.Unmarshal(body, &out)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}
	return nil, err
}

type DashboardResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      int    `json:"id"`
	URL     string `json:"url"`
	UID     string `json:"uid"`
}

type DashboardTargetDatasourceData struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type DashboardTargetSelectData struct {
	Type   string      `json:"type"`
	Params interface{} `json:"params"`
}

type DashboardTargetData struct {
	Datasource     DashboardTargetDatasourceData `json:"datasource"`
	Format         string                        `json:"format"`
	MetricColumn   string                        `json:"metricColumn"`
	RawQuery       bool                          `json:"rawQuery"`
	RawSQL         string                        `json:"rawSql"` // grafana need rawSql not rawSQL
	RefID          string                        `json:"refId"`
	Table          string                        `json:"table"`
	TimeColumn     string                        `json:"timeColumn"`
	TimeColumnType string                        `json:"timeColumnType"`
	Hide           bool                          `json:"hide"`
	Select         interface{}                   `json:"select"`
}

type GridPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type DashboardPanelData struct {
	Type          string      `json:"type"`
	Title         string      `json:"title"`
	PluginVersion string      `json:"pluginVersion"`
	Targets       interface{} `json:"targets"`
	GridPos       interface{} `json:"gridPos"`
}

type DashboardData struct {
	ID      string      `json:"id"`
	UID     string      `json:"uid"`
	Title   string      `json:"title"`
	Version int         `json:"version"`
	Panels  interface{} `json:"panels"`
	Refresh string      `json:"refresh"`
}

type BodyData struct {
	Dashboard DashboardData `json:"dashboard"`
}

type DashboardDetailResult struct {
	Dashboard interface{} `json:"dashboard"`
}

type DashboardDetail struct {
	Version string `json:"version"`
	ID      int    `json:"id"`
	UID     string `json:"uid"`
}
