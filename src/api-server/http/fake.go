package http

import (
	"log"
	"net"

	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/sqlite"
)

// Fake is a fake API Server HTTP server that sends the requests sequentially to the client.
type Server struct{}

// StartServer starts the gRPC server goroutine.
func (srv *Server) Start(cfg Config, pgClient *pg.Client) net.Addr {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Could not listen on ':0'")
	}

	cfg.Listen = lis

	go func() {
		StartHTTPService(cfg, pgClient)
	}()

	return lis.Addr()
}

// StartNewServer creates a Server and start the server.
func StartFakeNewServer(sqliteClient *sqlite.ORM, gLock *lock.Lock,
	waitCond *cond.Cond, pgClient *pg.Client, grafanaURL string,
) net.Addr {
	server := Server{}

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}
	nodeAgentDao := dao.NodeAgentDao{
		Client: sqliteClient,
	}

	moduleInstanceDao := dao.ModuleInstanceDao{
		Client: sqliteClient,
	}

	cfg := Config{
		GrafanaURL:      grafanaURL,
		GrafanaUserName: "admin",
		GrafanaUserPass: "admin",
		DatasourceName:  "TimescaleDB-Tricorder",
		DatasourceUID:   "timescaledb_tricorder",
		Module:          moduleDao,
		NodeAgent:       nodeAgentDao,
		ModuleInstance:  moduleInstanceDao,
		GLock:           gLock,
		WaitCond:        waitCond,
		Standalone:      false,
	}
	return server.Start(cfg, pgClient)
}
