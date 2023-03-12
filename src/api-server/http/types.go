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
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tricorder/src/api-server/http/dao"
	commonpb "github.com/tricorder/src/pb/module/common"
	"github.com/tricorder/src/pb/module/ebpf"
	"github.com/tricorder/src/pb/module/wasm"
	"github.com/tricorder/src/utils/pg"
)

// Common response type.
type HTTPResp struct {
	// Semantic and usage follow HTTP statues code convention.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
	Code int `json:"code"`

	// A human readable message explain the details of the status.
	Message string `json:"message"`
}

type CreateModuleReq struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	Wasm *wasm.Program `json:"wasm"`
	Ebpf *ebpf.Program `json:"ebpf"`
}

type CreateModuleResp struct {
	HTTPResp
}

type ListModuleReq struct {
	// These fields of the module record are returned to the client.
	// Empty list instructs server to return a default set of fields.
	// This allows client to control the size of the returned data to trade-off between responsiveness and completeness of
	// returned information.
	// TODO(yzhao): Change to string slice.
	Fields string
}

type ListModuleResp struct {
	HTTPResp
	Data []dao.ModuleGORM `json:"data"`
}

type ListAgentReq struct {
	// These fields of the agent record are returned to the client.
	// Empty list instructs server to return a default set of fields.
	// This allows client to control the size of the returned data to trade-off between responsiveness and completeness of
	// returned information.
	// TODO(yzhao): Change to string slice.
	Fields string
}

type ListAgentResp struct {
	HTTPResp
	Data []dao.NodeAgentGORM `json:"data"`
}

type DeployModuleResp struct {
	HTTPResp
	UID string `json:"uid"`
}

type UndeployModuleResp struct {
	HTTPResp
}

type DeleteModuleResp struct {
	HTTPResp
}

func checkQuery(c *gin.Context, key string) (string, error) {
	val, exist := c.GetQuery(key)
	if !exist {
		errMsg := fmt.Sprintf("key '%s' does not exist", key)
		c.JSON(http.StatusOK, gin.H{"code": "500", "message": errMsg})
		return "", fmt.Errorf(errMsg)
	}
	return val, nil
}

// DataFieldToPGColumn returns Column from a DataField protobuf message
func DataFieldToPGColumn(dataField *commonpb.DataField) (pg.Column, error) {
	return pg.Column{
		Name: dataField.Name,
		Type: dataField.Type,
	}, nil
}

// DataFieldsToPGColumns returns a slice of pg.Column for the input DataField slice.
func DataFieldsToPGColumns(dataFields []*commonpb.DataField) ([]pg.Column, error) {
	res := make([]pg.Column, 0, len(dataFields))
	for _, f := range dataFields {
		column, err := DataFieldToPGColumn(f)
		if err != nil {
			return nil, err
		}
		res = append(res, column)
	}
	return res, nil
}
