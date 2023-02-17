package grpc

import (
	"encoding/json"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/tricorder/src/api-server/dao"
	pb "github.com/tricorder/src/api-server/pb"
	modulepb "github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/pb/module/common"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	wasmpb "github.com/tricorder/src/pb/module/wasm"
	"github.com/tricorder/src/utils/channel"
)

type Deployer struct {
	Module dao.Module
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
		return err
	}
	log.Infof("Got input from client: %v", in)

	var eg errgroup.Group
	eg.Go(func() error {
		for {
			result, err := stream.Recv()
			if err != nil {
				log.Errorf("receive agent result error:%s", err.Error())
				return err
			}
			err = s.Module.UpdateStatusByID(result.ID, int(result.Status))
			if err != nil {
				log.Errorf("update code status error:%s", err.Error())
			}
		}
	})

	for {
		message := channel.ReceiveMessage()
		if message.Status != int(pb.DeploymentStatus_TO_BE_DEPLOYED) {
			continue
		}
		undeployList, _ := s.Module.ListCodeByStatus(int(pb.DeploymentStatus_TO_BE_DEPLOYED))
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
				ID:   code.ID,
				Name: code.Name,
				Module: &modulepb.Module{
					Ebpf: ebpf,
					Wasm: wasm,
				},
				Deploy: pb.DeployModuleReq_DEPLOY,
			}

			log.Infof("prepare to dispatch module: %v", &codeReq)
			err = stream.Send(&codeReq)
			if err != nil {
				// TODO(jian): The failure reason recorded in the err,
				// should be write into the sqlite database.instead of a logging message.
				log.Errorf("Deploy: [%s] failed: %s", code.Name, err.Error())

				err = s.Module.UpdateStatusByID(code.ID, int(pb.DeploymentStatus_DEPLOYMENT_FAILED))
				if err != nil {
					log.Errorf("update code status error:%s", err.Error())
				}
			} else {
				err = s.Module.UpdateStatusByID(code.ID, int(pb.DeploymentStatus_DEPLOYMENT_SUCCEEDED))
				if err != nil {
					log.Errorf("update code status error:%s", err.Error())
				}
			}
		}
	}
}
