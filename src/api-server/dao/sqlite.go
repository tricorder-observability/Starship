package dao

import (
	"fmt"

	"github.com/tricorder/src/utils/sqlite"

	log "github.com/sirupsen/logrus"
)

// InitSqlite prepares sqlite db file and setup the initial condition.
func InitSqlite(dbPath string) (*sqlite.ORM, error) {
	log.Infof("Opening SQLite database file at %s", dbPath)

	fullDbPath, err := sqlite.PrepareSqliteDbFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite datasource error %v", err)
	}
	engine, err := sqlite.NewORM(fullDbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite datasource error %v", err)
	}
	err = engine.CreateTable(&ModuleGORM{})
	if err != nil {
		return nil, fmt.Errorf("create code table error %v", err)
	}
	err = engine.CreateTable(&GrafanaAPIKeyGORM{})
	if err != nil {
		return nil, fmt.Errorf("create grafana_api table error %v", err)
	}
	return engine, nil
}
