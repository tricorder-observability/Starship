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
