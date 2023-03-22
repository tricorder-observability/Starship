package fake

import (
	"log"
	"net"

	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/wasm"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/sqlite"
)

// Fake is a fake API Server HTTP server that sends the requests sequentially to the client.
type Server struct{}

// StartServer starts the gRPC server goroutine.
func (srv *Server) Start(cfg http.Config, pgClient *pg.Client, wasiCompiler *wasm.WASICompiler) net.Addr {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Could not listen on ':0'")
	}

	cfg.Listen = lis

	go func() {
		err := http.StartHTTPService(cfg, pgClient, wasiCompiler)
		if err != nil {
			log.Fatalf("Failed to run HTTP Service, error: %v", err)
		}
	}()

	return lis.Addr()
}

// StartFakeNewServer creates a Server and start the server.
func StartFakeNewServer(
	sqliteClient *sqlite.ORM, gLock *lock.Lock,
	waitCond *cond.Cond, pgClient *pg.Client, grafanaURL string,
	wasiSDKPath string, wasiStarshipIncludePath string,
	wasiBuildTmpPath string,
) net.Addr {
	server := Server{}

	dao := dao.NewDao(sqliteClient)

	cfg := http.Config{
		GrafanaURL:      grafanaURL,
		GrafanaUserName: "admin",
		GrafanaUserPass: "admin",
		DatasourceName:  "TimescaleDB-Tricorder",
		DatasourceUID:   "timescaledb_tricorder",
		Module:          dao.Module,
		NodeAgent:       dao.NodeAgent,
		ModuleInstance:  dao.ModuleInstance,
		GLock:           gLock,
		WaitCond:        waitCond,
		Standalone:      false,
	}
	wasiCompiler := wasm.NewWASICompiler(wasiSDKPath, wasiStarshipIncludePath, wasiBuildTmpPath)
	return server.Start(cfg, pgClient, wasiCompiler)
}
