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
	commonpb "github.com/tricorder/src/pb/module/common"
	"github.com/tricorder/src/utils/pg"
)

// CreateModuleResponse binds to the HTTP response sent to the management Web UI.
type CreateModuleResponse struct {
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
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
