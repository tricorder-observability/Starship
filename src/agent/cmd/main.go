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

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/agent/deployer"
	proc_info "github.com/tricorder/src/agent/proc-info"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"

	linux_headers "github.com/tricorder/src/agent/ebpf/bcc/linux-headers"
	"github.com/tricorder/src/agent/ebpf/bcc/utils"
)

var (
	// For compatiability, --module_deployer_address not rename for now
	apiServerAddr = flag.String(
		"module_deployer_address",
		"localhost:50051",
		"The grpc address of API Server's ModuleDeployer and ProcessCollector service",
	)
	modulePGURL = flag.String("pg_url", "postgresql://postgres:password@localhost", "The URL to PostgreSQL instance")
	// The default value is incompatible with the container environment, which mounts the host's `/` to `/host` inside
	// the container. This is for easier testing during development, which is not inside a container.
	hostSysRootPath = flag.String("host_sys_root_path", "/sys", "The path to the host's /sys file system that "+
		"can be accessed by agent, this is mounted by Kubernetes. Tricorder reads cgroup and BPF probes from files "+
		"under this directory")
)

func main() {
	flag.Parse()

	cfg, err := newConfig()
	if err != nil {
		log.Fatalf("Failed to create config, error: %v", err)
	}

	if err := utils.CleanTricorderProbes(*hostSysRootPath); err != nil {
		log.Warnf("Failed to cleanup previously-deployed dangling probes, error: %v", err)
	}

	if err := linux_headers.Init(); err != nil {
		log.Errorf("Failed to initialize Linux headers for bcc, error: %v", err)
	}

	deployer := deployer.New(*apiServerAddr, cfg.nodeName, cfg.podID)
	for {
		communicateWithNode(cfg.nodeName, deployer)
	}

	deployer.Stop()
	log.Infof("Hello deployer\n")
}

func communicateWithNode(nodeName string, deployer *deployer.Deployer) error {
	err := retry.ExpBackOffWithLimit(func() error {
		return deployer.ConnectToAPIServer()
	})
	if err != nil {
		log.Errorf("Failed to connect to API server, error: %v", err)
		return err
	}

	pgClient := pg.NewClient(*modulePGURL)
	err = retry.ExpBackOffWithLimit(pgClient.Connect)
	if err != nil {
		log.Errorf("Failed to connect to database at '%s', error: %v", *modulePGURL, err)
		return err
	}
	deployer.PGClient = pgClient

	collector := proc_info.NewCollector(*hostSysRootPath, *apiServerAddr, nodeName)
	if err := collector.StartProcInfoReport(); err != nil {
		log.Errorf("Failed to ReportProcess, error: %v", err)
		return err
	}

	return deployer.StartModuleDeployLoop()
}
