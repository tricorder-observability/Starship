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

package deployer

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/tricorder/src/agent/ebpf/bcc/linux_headers"
	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/pb/module/common"
	ebpf "github.com/tricorder/src/pb/module/ebpf"
	"github.com/tricorder/src/pb/module/wasm"
	testutils "github.com/tricorder/src/testing/bazel"
)

// TODO(yzhao): Rewrite this test to use real API server and agent.

type GRPCDeployerHandler struct {
	// mockV2QueryService *querysvc.QueryService
}

const code string = `
#include <linux/ptrace.h>

BPF_PERF_OUTPUT(events);

int syscall__probe_entry_read(struct pt_regs* ctx, int fd, char* buf, size_t count) {
	bpf_trace_printk("syscall__probe_entry_read\n");
  return 0;
}
int syscall__probe_return_read(struct pt_regs* ctx) {
	bpf_trace_printk("syscall__probe_ret_read\n");
  return 0;
}
`

func mockDeployReqs() ([]*pb.DeployModuleReq, error) {
	wasmRelPath := "modules/sample_json/copy_input_to_output.wasm"
	wasmBinaryCode, err := testutils.ReadTestBinFile(wasmRelPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read wasm file %s, error: %v", wasmRelPath, err)
	}
	return []*pb.DeployModuleReq{
		{},
		{
			ID:   "mock_empty_ebpf_code-1",
			Name: "mock_empty_ebpf_code",
		},
		{
			ID:             "mock_test_deploy_module_req-1",
			Name:           "mock_test_deploy_module_req",
			LifeTimeSecond: 0,
			Module: &module.Module{
				Ebpf: &ebpf.Program{
					Fmt:            common.Format_TEXT,
					Lang:           common.Lang_C,
					Code:           code,
					PerfBufferName: "events",
					Probes: []*ebpf.ProbeSpec{
						{
							Type:   ebpf.ProbeSpec_SYSCALL_PROBE,
							Target: "read",
							Entry:  "syscall__probe_entry_read",
							Return: "syscall__probe_return_read",
						},
					},
				},
				Wasm: &wasm.Program{
					Fmt:    common.Format_BINARY,
					Code:   wasmBinaryCode,
					FnName: "copy_input_to_output",
					OutputSchema: &common.Schema{
						Name: "data",
						Fields: []*common.DataField{
							{
								Name: "data",
								Type: common.DataField_JSONB,
							},
						},
					},
				},
			},

			Deploy: pb.DeployModuleReq_DEPLOY,
		},
		{
			ID:     "mock_test_deploy_module_req-1",
			Name:   "mock_test_deploy_module_req",
			Deploy: pb.DeployModuleReq_UNDEPLOY,
		},
	}, nil
}

func (g *GRPCDeployerHandler) DeployModule(stream pb.ModuleDeployer_DeployModuleServer) error {
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Printf("Got input from client: %v", in)

	reqs, err := mockDeployReqs()
	if err != nil {
		return err
	}
	for _, req := range reqs {
		err = stream.Send(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func newMockGRPCServer(t *testing.T) (*grpc.Server, net.Addr) {
	lis, _ := net.Listen("tcp", ":0")
	grpcServer := grpc.NewServer()

	grpcHandler := &GRPCDeployerHandler{}
	pb.RegisterModuleDeployerServer(grpcServer, grpcHandler)

	go func() {
		err := grpcServer.Serve(lis)
		require.NoError(t, err)
	}()

	return grpcServer, lis.Addr()
}

func TestDeployAndRun(t *testing.T) {
	assert := assert.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())

	_, addr := newMockGRPCServer(t)

	d := New(addr.String(), "node_name", "pid_id")

	d.apiServerAddr = addr.String()
	err := d.ConnectToAPIServer()
	require.NoError(t, err)

	err = d.InitModuleDeployLink()
	require.NoError(t, err)

	err = d.StartModuleDeployLoop()
	require.NoError(t, err)

	// this module has been deploy and then undeploy
	assert.Nil(d.idDeployMap["mock_test_deploy_module_req-1"])

	d.Stop()
}
