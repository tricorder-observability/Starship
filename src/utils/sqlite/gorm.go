package sqlite

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A Object–relational mapping (https://en.wikipedia.org/wiki/Object-relational_mapping)
// type for intermediating between Golang types and Sqlite DB.
type ORM struct {
	// The underlying engine that provides all APIs.
	Engine *gorm.DB
}

// NewROM Returns a new ORM object.
func NewORM(dbfile string) (*ORM, error) {
	client := new(ORM)
	engine, err := gorm.Open(sqlite.Open(dbfile),
		// See https://gorm.io/docs/gorm_config.html for detailed configurations.
		&gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to database at '%s', error: %v", dbfile, err)
	}
	client.Engine = engine
	return client, nil
}

// CreateTable wraps gorm.DB.AutoMigrate()
func (g *ORM) CreateTable(schema interface{}) error {
	return g.Engine.AutoMigrate(schema)
}
