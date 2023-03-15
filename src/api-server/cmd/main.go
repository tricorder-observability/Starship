// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	sg "github.com/tricorder/src/api-server/grpc"
	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/http/docs"
	"github.com/tricorder/src/api-server/meta"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/errors"
	grpcutils "github.com/tricorder/src/utils/grpc"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/log"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"
)

var (
	testOnlyHost = flag.String("test_only_host", "localhost",
		"The host address used for displaying swagger, this allows Swagger UI to connect to the running server.")

	standalone = flag.Bool("standalone", false, "If true, API Server can be started without dependent services")

	// Management Web UI requires to connect to Postgres, Grafana, this allows us to disable this service in tests.
	enableMgmtUI = flag.Bool("enable_mgmt_ui", true, "If true, start management Web UI")
	// This allows us to disable this service in tests.
	enableGRPC = flag.Bool("enable_grpc", true, "If true, start the gRPC server for managing agents")

	// Metadata service requires to connect to Postgres, this allows us to disable this service in tests.
	enableMetadataService = flag.Bool(
		"enable_metadata_service",
		true,
		"If true, start collecting metadata from K8s API Server and write to Postgres",
	)
	// For compatiability, module_deployer_port not rename for now
	agentServicePort = flag.Int("module_deployer_port", 50051, "The port to which the ModuleDeployer service listens")

	modulePGURL = flag.String("pg_url", "postgresql://postgres:password@localhost", "The URL to PostgreSQL instance")

	mgmtUIPort      = flag.Int("mgmt_ui_port", 8080, "The port to which the management Web UI listens")
	moduleDBDirPath = flag.String(
		"module_db_dir_path",
		"src/api-server/http/",
		"The dir path to the SQLite database file",
	)
	moduleGrafanaURL = flag.String("grafana_url", "http://localhost:3000", "The URL to PostgreSQL instance")

	// These 2 flags must be the same as Grafana configuration in helm-charts's charts/starship/values.yaml
	moduleGrafanaUserName = flag.String("grafana_user_name", "admin",
		"Grafana username, must be consistent with Grafana installation config")
	moduleGrafanaUserPassword = flag.String("grafana_user_pass", "tricorder",
		"Grafana password, must be consistent with Grafana installation config")

	moduleDatasourceName = flag.String(
		"grafana_ds_name",
		"TimescaleDB-Tricorder",
		"The name of datasource in grafana to postgres database",
	)
	moduleDatasourceUID = flag.String(
		"grafana_ds_uid",
		"timescaledb_tricorder",
		"The uid of datasource in grafana to postgres database",
	)
)

func setupSwaggerInfo() {
	docs.SwaggerInfo.Title = "API Server"
	docs.SwaggerInfo.Description = "API Server http api document."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", *testOnlyHost, *mgmtUIPort)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func main() {
	flag.Parse()

	setupSwaggerInfo()

	log.Infof("Creating Postgresql client at %s", *modulePGURL)
	pgClient := pg.NewClient(*modulePGURL)
	// global K8s clientset, shared with lots of informers
	var clientset kubernetes.Interface

	err := retry.ExpBackOffWithLimit(func() error {
		return pgClient.Connect()
	})
	if err != nil {
		log.Fatalf("Failed to initialize a client to Postgresql database error at %s, error: %v", *modulePGURL, err)
	}

	sqliteClient, _ := dao.InitSqlite(*moduleDBDirPath)
	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}
	nodeAgentDao := dao.NodeAgentDao{
		Client: sqliteClient,
	}

	moduleInstanceDao := dao.ModuleInstanceDao{
		Client: sqliteClient,
	}

	waitCond := cond.NewCond()
	gLock := lock.NewLock()

	if *enableMetadataService {
		err = retry.ExpBackOffWithLimit(func() error {
			clientset, err = kubernetes.NewForConfig(ctrl.GetConfigOrDie())
			return err
		})
		if err != nil {
			log.Fatalf("while starting resource watching service, failed to create K8s client, error: %v", err)
		}
	}

	// Launch all long-running server goroutines.
	var srvErrGroup errgroup.Group

	if *enableGRPC {
		srvErrGroup.Go(func() error {
			f, err := grpcutils.NewServerFixture(*agentServicePort)
			if err != nil {
				return errors.Wrap("starting gRPC server", "create server fixture", err)
			}
			sg.RegisterModuleDeployerServer(f, sqliteClient, gLock, waitCond)
			if *enableMetadataService {
				sg.RegisterProcessCollectorServer(f, clientset, pgClient)
			}
			err = f.Serve()
			if err != nil {
				return errors.Wrap("starting gRPC server", "serve", err)
			}
			return nil
		})
	}

	if *enableMgmtUI {
		srvErrGroup.Go(func() error {
			const tcp = "tcp"
			addrStr := fmt.Sprintf(":%d", *mgmtUIPort)
			listener, err := net.Listen(tcp, addrStr)
			if err != nil {
				return errors.Wrap("starting http server", "listen", err)
			}
			config := http.Config{
				Listen:          listener,
				GrafanaURL:      *moduleGrafanaURL,
				GrafanaUserName: *moduleGrafanaUserName,
				GrafanaUserPass: *moduleGrafanaUserPassword,
				DatasourceName:  *moduleDatasourceName,
				DatasourceUID:   *moduleDatasourceUID,
				Module:          moduleDao,
				NodeAgent:       nodeAgentDao,
				ModuleInstance:  moduleInstanceDao,
				WaitCond:        waitCond,
				GLock:           gLock,
				Standalone:      *standalone,
			}

			http.StartHTTPService(config, pgClient)
			return nil
		})
	}

	if *enableMetadataService {
		srvErrGroup.Go(func() error {
			err := meta.StartWatchingResources(clientset, pgClient, &nodeAgentDao, waitCond)
			if err != nil {
				log.Fatalf("Could not start metadata service, error: %v", err)
			}
			return nil
		})
	}

	log.Infof("API server has started ...")

	srvErr := srvErrGroup.Wait()
	if srvErr != nil {
		log.Fatalf("Server goroutines failed, error: %v", srvErr)
	}
}
