package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/api-server/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	testutils "github.com/tricorder/src/testing/bazel"
	grafanaTest "github.com/tricorder/src/testing/grafana"
)

// Tests that we can initialize API token on Grafana.
func TestInitGrafanaAPIToken(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cleanerFn, grafanaURL, err := grafanaTest.LaunchContainer()
	require.Nil(err)
	defer func() {
		assert.Nil(cleanerFn())
	}()

	grafana.InitGrafanaConfig(grafanaURL, "admin", "admin")
	sqliteClient, err := dao.InitSqlite(testutils.GetTmpFile())
	assert.Nil(err)

	grafanaManager := GrafanaManagement{
		GrafanaAPIKey: dao.GrafanaAPIKey{
			Client: sqliteClient,
		},
	}
	assert.Nil(grafanaManager.InitGrafanaAPIToken())
}
