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
