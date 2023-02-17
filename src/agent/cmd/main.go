package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	_ "github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/agent/deployer"
	"github.com/tricorder/src/agent/proc_info"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"

	"github.com/tricorder/src/agent/ebpf/bcc/linux_headers"
	"github.com/tricorder/src/agent/ebpf/bcc/utils"
)

var (
	// For compatiability, --module_deployer_address not rename for now
	apiServerAddr = flag.String(
		"module_deployer_address",
		"localhost:50051",
		"The address of API Server's ModuleDeployer and ProcessCollector service",
	)
	modulePGURL     = flag.String("pg_url", "postgresql://postgres:password@localhost", "The URL to PostgreSQL instance")
	hostSysRootPath = flag.String("host_sys_root_path", "/host/sys", "The path to the host's /sys file system that "+
		"can be accessed by agent, this is mounted by Kubernetes. Tricorder reads cgroup and BPF probes from files "+
		"under this directory")
)

func main() {
	flag.Parse()

	if err := utils.CleanTricorderProbes(*hostSysRootPath); err != nil {
		log.Warnf("Failed to cleanup previously-deployed dangling probes, error: %v", err)
	}

	if err := linux_headers.Init(); err != nil {
		log.Errorf("Failed to initialize Linux headers for bcc, error: %v", err)
	}

	var deployer deployer.Deployer

	err := retry.ExpBackOffWithLimit(func() error {
		return deployer.ConnectToAPIServer(*apiServerAddr)
	})
	if err != nil {
		log.Fatalf("Failed to connect to API server, error: %v", err)
	}

	pgClient := pg.NewClient(*modulePGURL)
	err = retry.ExpBackOffWithLimit(pgClient.Connect)
	if err != nil {
		log.Fatalf("Failed to connect to database at '%s', error: %v", *modulePGURL, err)
	}
	deployer.PGClient = pgClient

	err = retry.ExpBackOffWithLimit(deployer.InitModuleDeployLink)
	if err != nil {
		log.Fatalf("Failed to establish stream connection to module deploy service, error: %v", err)
	}

	collector := proc_info.NewCollector(*hostSysRootPath, *apiServerAddr)
	if err := collector.StartProcInfoReport(); err != nil {
		log.Errorf("Failed to ReportProcess, error: %v", err)
	}

	err = deployer.StartModuleDeployLoop()
	if err != nil {
		log.Fatalf("Failed to start deployment loop, error: %v", err)
	}

	deployer.Stop()
	log.Infof("Hello deployer\n")
}
