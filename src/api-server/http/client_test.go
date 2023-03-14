package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/http/grafana"
	testutils "github.com/tricorder/src/testing/bazel"
	grafanatest "github.com/tricorder/src/testing/grafana"
	pgclienttest "github.com/tricorder/src/testing/pg"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
)

func TestListAgents(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	grafana.InitGrafanaConfig(grafanaURL, "admin", "admin")

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	FakeHTTPServer := StartFakeNewServer(sqliteClient, gLock, waitCond, pgClient, grafanaURL)

	// test list agent
	client := NewClient("http://" + FakeHTTPServer.String())
	res, err := client.ListAgents(nil)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Equal(0, len(res.Data))

	nodeAgentDao := dao.NodeAgentDao{
		Client: sqliteClient,
	}
	newAgent := dao.NodeAgentGORM{
		AgentID:    "agent_test_id",
		NodeName:   "agent_test_node",
		AgentPodID: "agent_test_pod_id",
	}
	err = nodeAgentDao.SaveAgent(&newAgent)
	require.NoError(err)

	res, err = client.ListAgents(nil)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Equal(1, len(res.Data))
	assert.Equal("agent_test_id", res.Data[0].AgentID)
}
