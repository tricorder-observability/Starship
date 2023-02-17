package proc_info

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	p "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/retry"
)

// Collector is responsible for grab process info with containers/pods.
type Collector struct {
	// The base path to the host's /sys file system (the Kubernetes node).
	// The association between local process ID and container ID is established by looking up process ID list of cgroup
	// under /sys file system. And we mount the host's /sys file system to agent's container to facilitate lookup.
	hostSysRootPath string

	// The address to the API server.
	apiServerAddr string

	// Connects to API server's process info collector server, and reports process information.
	procCollectorClient pb.ProcessCollectorClient
}

func NewCollector(hostSysRootPath, apiServerAddr string) *Collector {
	return &Collector{
		hostSysRootPath: hostSysRootPath,
		apiServerAddr:   apiServerAddr,
	}
}

func (c *Collector) connect() error {
	grpcConn, err := grpc.Dial(c.apiServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to API server at '%s', error: %v", c.apiServerAddr, err)
	}
	c.procCollectorClient = pb.NewProcessCollectorClient(grpcConn)
	return nil
}

func (c *Collector) StartProcInfoReport() error {
	err := retry.ExpBackOffWithLimit(c.connect)
	if err != nil {
		log.Fatalf("Failed to connect to API server, error: %v", err)
	}

	stream, err := c.procCollectorClient.ReportProcess(context.Background())
	if err != nil {
		return err
	}

	if err = stream.Send(&pb.ProcessWrapper{Msg: &pb.ProcessWrapper_NodeName{NodeName: GetNodeName()}}); err != nil {
		log.Errorf("stream.Send error: %v", err)
	}

	go func() {
		for {
			containerInfo, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("while report process info, gRPC stream to API server broke, error: %v", err)
			}
			processInfo, err := grabProcessInfo(c.hostSysRootPath+"/fs/cgroup", containerInfo)
			if err != nil {
				log.Errorf("while collecting process info for container '%v', "+
					"failed to grab process info, error: %v", containerInfo, err)
				continue
			}

			if err = stream.Send(&pb.ProcessWrapper{Msg: &pb.ProcessWrapper_Process{Process: processInfo}}); err != nil {
				log.Errorf("stream.Send error: %v", err)
			}
		}
	}()

	return nil
}

// GetNodeName returns value injected by downwardAPI
// Inject outer-scope hostname into container, so the agent can use this to filter out updates not relevant to this node
// from the K8s API server.
// env:
// - name: NODE_NAME
//   valueFrom:
//	   fieldRef:
//	     fieldPath: spec.nodeName
func GetNodeName() string {
	return os.Getenv("NODE_NAME")
}

func getProcCreateTime(pid int32) (int64, error) {
	proc, err := p.NewProcess(pid)
	if err != nil {
		return -1, err
	}
	return proc.CreateTime()
}

func grabProcessInfo(basePath string, ci *pb.ContainerInfo) (*pb.ProcessInfo, error) {
	procList := []*pb.Process{}
	found := false
	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Container's ID in the format '<type>://<container_id>'.
			if strings.Contains(info.Name(), strings.Split(ci.Id, "://")[1]) && file.Exists(path+"/cgroup.procs") {
				procList, err = readPIDsAndCreateTime(path + "/cgroup.procs")
				if err != nil {
					return err
				}
				found = true
			}
			return nil
		})

	if !found {
		return nil, fmt.Errorf("not find container %s[%s] of pod %s[%s] pid file", ci.Name, ci.Id, ci.PodName, ci.PodUid)
	}

	return &pb.ProcessInfo{ProcList: procList, Container: ci}, err
}

func readPIDsAndCreateTime(fullpath string) ([]*pb.Process, error) {
	procList := []*pb.Process{}
	lines, err := file.ReadLines(fullpath)
	if err != nil {
		return nil, fmt.Errorf("while grabbing process info, failed to read proc.status file, error: %v", err)
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		pid, err := strconv.ParseInt(line, 10 /*base*/, 32 /*bitSize*/)
		if err != nil {
			return nil, fmt.Errorf("while grabbing process info, failed to parse PID, error: %v", err)
		}
		createTime, err := getProcCreateTime(int32(pid))
		if err != nil {
			return nil, fmt.Errorf("while grabbing process info, failed to get creation time, error: %v", err)
		}
		procList = append(procList, &pb.Process{Id: int32(pid), CreateTime: createTime})
	}
	return procList, nil
}
