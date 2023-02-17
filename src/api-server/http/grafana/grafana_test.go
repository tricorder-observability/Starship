package grafana

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	grafanaTest "github.com/tricorder/src/testing/grafana"
)

// Tests that auth token can be created on Grafana
func TestAuthToken(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cleanerFn, grafanaURL, err := grafanaTest.LaunchContainer()
	require.Nil(err)
	defer func() {
		// Have to write in this form
		// defer assert.Nil(cleanerFn())
		// causes cleanerFn() be invoked immediately, not really deferred.
		assert.Nil(cleanerFn())
	}()

	log.Println("grafana url:" + grafanaURL)
	InitGrafanaConfig(grafanaURL, "admin", "admin")
	authToken := NewAuthToken()
	require.NotNil(authToken)

	token, err := authToken.GetToken("/api/dashboards/db")
	assert.Nil(err)

	log.Println("create grafana token success: ", token)
	assert.Nil(createDashboard(token.Key))

	datasourceToken, err := authToken.GetToken("/api/datasources")
	require.Nil(err)
	assert.Nil(createDatasource(datasourceToken.Key))
}

func createDashboard(token string) error {
	dashboard := NewDashboard()
	if dashboard == nil {
		return fmt.Errorf("failed to create dashboard")
	}
	result, err := dashboard.CreateDashboard(token, "APIServer1", "uid")
	if err != nil {
		return fmt.Errorf("create grafana dashboard error:%v", err)
	}
	if result.Status != "success" {
		return fmt.Errorf("create grafana dashboard error:%s: %s", result.Status, result.Message)
	}
	return nil
}

func createDatasource(token string) error {
	ds := NewDatasource()
	const name = "MySQLTEST"

	if ds == nil {
		return fmt.Errorf("failed to create datasource")
	}

	_, err := ds.CreateDatasource(token, name,
		"localhost:5432", "postgres", "123456", "test")
	if err != nil {
		return fmt.Errorf("create grafana datasource error:%v", err)
	}
	return nil
}
