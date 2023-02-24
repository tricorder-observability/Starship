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

package grpc

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/sync/errgroup"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/api-server/dao"
	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/api-server/utils/channel"
	modulepb "github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/pb/module/common"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	wasmpb "github.com/tricorder/src/pb/module/wasm"
)

// Manages the deployment of eBPF+WASM modules
type Deployer struct {
	// The DAO object that proxies with SQLite for writing and reading the serialized data.
	Module dao.Module

	// The list of agents connected with this Deployer.
	//
	// Each agent and this Deployer maintains a gRPC streaming channel with DeployModuleReq & DeployModuleResp
	// flow back-and-forth.
	agents []*pb.Agent
}

// DeployModule implements the only RPC of the ModuleDeployer service.
// It continuously sends deployment request to the connected agent (as client).
func (s *Deployer) DeployModule(stream pb.ModuleDeployer_DeployModuleServer) error {
	// The first message is sent from the client, but the remaining loops are driven by the server.
	// The server will send deploy module request for this client to work on.
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return errors.Wrap("handling DeployModule", "receive message", err)
	}

	log.Infof("Agent '%s' connected, starting module management loop ...", in.Agent.Id)

	s.agents = append(s.agents, in.Agent)

	var eg errgroup.Group

	// Create a goroutine to check the response from the connected agent.
	eg.Go(func() error {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				log.Warnf("Agent closed connection, this should only happens during testing; stopping ...")
				return nil
			}
			if err != nil {
				return errors.Wrap("handling DeployModule request", "receive mssage", err)
			}
			// Need to be able to correctly account which nodes are deployed, and which are
			// not deployed. TODO(yzhao & zhihui): Chat about the design.
			err = s.Module.UpdateStatusByID(result.ModuleId, int(result.State))
			if err != nil {
				log.Errorf("update code status error:%s", err.Error())
			}
		}
	})

	for {
		// TODO(yzhao): This should be moved into gRPC side, not in utils.
		message := channel.ReceiveMessage()
		if message.Status != int(pb.DeploymentState_TO_BE_DEPLOYED) {
			continue
		}
		undeployList, _ := s.Module.ListCodeByStatus(int(pb.DeploymentState_TO_BE_DEPLOYED))
		for _, code := range undeployList {
			var probeSpecs []*ebpfpb.ProbeSpec
			if len(code.EbpfProbes) > 0 {
				err = json.Unmarshal([]byte(code.EbpfProbes), &probeSpecs)
				if err != nil {
					return fmt.Errorf("while deploying module, failed to unmarshal ebpf probes, error: %v", err)
				}
			}

			ebpf := &ebpfpb.Program{
				Fmt:            common.Format(code.EbpfFmt),
				Lang:           common.Lang(code.EbpfLang),
				Code:           code.Ebpf,
				PerfBufferName: code.EbpfPerfBufferName,
				Probes:         probeSpecs,
			}

			var fields []*common.DataField
			if len(code.SchemaAttr) > 0 {
				err = json.Unmarshal([]byte(code.SchemaAttr), &fields)
				if err != nil {
					return fmt.Errorf("while deploying module, failed to unmarshal data fields, error: %v", err)
				}
			}

			wasm := &wasmpb.Program{
				Fmt:    common.Format(code.WasmFmt),
				Lang:   common.Lang(code.WasmLang),
				FnName: code.Fn,
				OutputSchema: &common.Schema{
					Name:   code.SchemaName,
					Fields: fields,
				},
				Code: code.Wasm,
			}

			codeReq := pb.DeployModuleReq{
				ModuleId: code.ID,
				Module: &modulepb.Module{
					Ebpf: ebpf,
					Wasm: wasm,
				},
				Deploy: pb.DeployModuleReq_DEPLOY,
			}

			err = stream.Send(&codeReq)
			if err != nil {
				// TODO(jian): The failure reason recorded in the err,
				// should be write into the sqlite database.instead of a logging message.
				log.Errorf("Deploy: [%s] failed: %s", code.Name, err.Error())

				err = s.Module.UpdateStatusByID(code.ID, int(pb.DeploymentState_DEPLOYMENT_FAILED))
				if err != nil {
					log.Errorf("update code status error:%s", err.Error())
				}
			} else {
				err = s.Module.UpdateStatusByID(code.ID, int(pb.DeploymentState_DEPLOYMENT_SUCCEEDED))
				if err != nil {
					log.Errorf("update code status error:%s", err.Error())
				}
			}
		}
	}
}
