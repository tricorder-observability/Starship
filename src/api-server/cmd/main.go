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

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/api-server/dao"
	sg "github.com/tricorder/src/api-server/grpc"
	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/meta"
	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"
)

var (
	// Management Web UI requires to connect to Postgres, Grafana, this allows us to disable this service in tests.
	enableMgmtUI = flag.Bool("enable_mgmt_ui", true, "If true, start management Web UI")

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

func main() {
	flag.Parse()

	// Log the line number.
	log.SetReportCaller(true)

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
	codeDao := dao.Module{
		Client: sqliteClient,
	}
	grafanaAPIDao := dao.GrafanaAPIKey{
		Client: sqliteClient,
	}

	if *enableMetadataService {
		err = retry.ExpBackOffWithLimit(func() error {
			clientset, err = kubernetes.NewForConfig(ctrl.GetConfigOrDie())
			return err
		})
		if err != nil {
			log.Fatalf("while starting resource watching service, failed to create K8s client, error: %v", err)
		}
	}

	var eg errgroup.Group
	eg.Go(func() error {
		// Agent service server side, including module deployer and process collector service
		err := startAgentServerSide(*agentServicePort, codeDao, pgClient, clientset)
		if err != nil {
			log.Fatalf("Could not start server, error: %v", err)
		}

		return nil
	})

	if *enableMgmtUI {
		eg.Go(func() error {
			config := http.Config{
				Port:            *mgmtUIPort,
				GrafanaURL:      *moduleGrafanaURL,
				GrafanaUserName: *moduleGrafanaUserName,
				GrafanaUserPass: *moduleGrafanaUserPassword,
				DatasourceName:  *moduleDatasourceName,
				DatasourceUID:   *moduleDatasourceUID,
				Module:          codeDao,
				GrafanaAPIKey:   grafanaAPIDao,
			}

			http.StartHTTPService(config, pgClient)
			return nil
		})
	}

	if *enableMetadataService {
		eg.Go(func() error {
			if err := meta.StartWatchingResources(clientset, pgClient); err != nil {
				log.Errorf("Could not start metadata service, error: %v", err)
			}
			return nil
		})
	}

	log.Infof("API server has started ...")
	_ = eg.Wait()
}

func startAgentServerSide(port int, c dao.Module, pgClient *pg.Client, clientset kubernetes.Interface) error {
	addr := fmt.Sprintf(":%d", port)
	log.Infof("Starting gRPC server at %s", addr)

	grpcLis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "Failed to start listening on gRPC address %s", addr)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterModuleDeployerServer(grpcServer, &sg.Deployer{Module: c})
	if *enableMetadataService {
		pb.RegisterProcessCollectorServer(grpcServer, sg.NewPIDCollector(clientset, pgClient))
	}
	err = grpcServer.Serve(grpcLis)

	if err != nil {
		return errors.Wrap(err, "start grpc server")
	}
	return nil
}
