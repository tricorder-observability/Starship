package pg

import (
	"fmt"
	"strings"

	commonpb "github.com/tricorder/src/pb/module/common"
)

// https://www.w3schools.com/sql/sql_constraints.asp
const (
	NOT_NULL     = "NOT NULL"
	UNIQUE       = "UNIQUE"
	PRIMARY_KEY  = "PRIMARY KEY"
	FOREIGN_KEY  = "FOREIGN KEY"
	CHECK        = "CHECK"
	DEFAULT      = "DEFAULT"
	CREATE_INDEX = "CREATE INDEX"
)

var DataTypeConstraints = map[string]bool{
	NOT_NULL:     true,
	UNIQUE:       true,
	PRIMARY_KEY:  true,
	FOREIGN_KEY:  true,
	CHECK:        true,
	DEFAULT:      true,
	CREATE_INDEX: true,
}

// Convenience consts to save the needs of importing pb field and long name.
const (
	BOOL    = commonpb.DataField_BOOL
	DATE    = commonpb.DataField_DATE
	INT     = commonpb.DataField_INT
	INTEGER = commonpb.DataField_INTEGER
	JSON    = commonpb.DataField_JSON
	JSONB   = commonpb.DataField_JSONB
	TEXT    = commonpb.DataField_TEXT
)

// Column describes a column of a data table in a database.
type Column struct {
	Name       string
	Type       commonpb.DataField_Type
	Constraint string
}

// Returns a string that defines this column in a SQL expression.
func DefineColumn(c Column) (string, error) {
	if _, ok := DataTypeConstraints[c.Constraint]; len(c.Constraint) != 0 && !ok {
		return "", fmt.Errorf("while defining column '%s', constraint '%s' is not supported", c.Name, c.Constraint)
	}
	typeName, ok := commonpb.DataField_Type_name[int32(c.Type)]
	if !ok {
		return "", fmt.Errorf("while defining column '%s', data type '%s' is not supported", c.Name, c.Type)
	}
	if len(c.Constraint) == 0 {
		return strings.Join([]string{c.Name, typeName}, " "), nil
	}
	return strings.Join([]string{c.Name, typeName, c.Constraint}, " "), nil
}
