package dao

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
)

// test GrafanaAPIKey fun
// test SaveGrafanaAPI fun and check save result
// test QueryByID and QueryByAPIKey and check query result
func TestGrafanaAPIDAO(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	defer func() {
		assert.Nil(os.RemoveAll(dirPath))
	}()
	sqliteClient, _ := InitSqlite(dirPath)

	grafanaAPIDao := GrafanaAPIKey{
		Client: sqliteClient,
	}

	grafanaAPIKey := &GrafanaAPIKeyGORM{
		Name:       "testAPIName",
		APIKEY:     "APIKey",
		AuthValue:  "AuthTokenValue",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := grafanaAPIDao.SaveGrafanaAPI(grafanaAPIKey)
	if err != nil {
		t.Errorf("save grafana api err %v", err)
	}

	queryResult, err := grafanaAPIDao.QueryByID(grafanaAPIKey.ID)
	if err != nil {
		t.Errorf("query grafana api by id err %v", err)
	}
	if grafanaAPIKey.Name != queryResult.Name {
		t.Errorf("save grafana api by id error: not found name == %v", grafanaAPIKey.Name)
	}

	queryResult, err = grafanaAPIDao.QueryByAPIKey(grafanaAPIKey.APIKEY)
	if err != nil {
		t.Errorf("query grafana api by id err %v", err)
	}
	if grafanaAPIKey.APIKEY != queryResult.APIKEY {
		t.Errorf("save grafana api by api key error: not found apiKey == %v", grafanaAPIKey.APIKEY)
	}
}
