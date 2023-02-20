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

package pg

import "github.com/tricorder/src/pb/module/common"

// DataTable describes a data table for writing data.
type Schema struct {
	// Name of the data table
	Name string

	// All fields
	Columns []Column
}

// jsonbTableSchemaTmpl defines the data table columns, only requires to define name.
var jsonbTableSchemaTmpl = Schema{
	Columns: []Column{
		{
			Name: "data",
			Type: JSONB,
		},
	},
}

// GetJSONBTableSchema returns table schema with JSON data type and set the name.
func GetJSONBTableSchema(tableName string) *Schema {
	schema := jsonbTableSchemaTmpl
	schema.Name = tableName
	return &schema
}

// Returns a Schema from a protobuf with the same semantics.
func SchemaFromPB(pbSchema *common.Schema) *Schema {
	columns := make([]Column, 0, len(pbSchema.Fields))
	for _, field := range pbSchema.Fields {
		columns = append(columns, Column{Name: field.Name, Type: field.Type})
	}
	return &Schema{
		Name:    pbSchema.Name,
		Columns: columns,
	}
}
