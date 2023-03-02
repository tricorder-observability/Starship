package fake

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tricorder/src/api-server/pb"
)

// Fake is a fake API Server gRPC server that sends the requests sequentially to the client.
type Server struct {
	Reqs []*pb.DeployModuleReq
}

// Implements the API Server's gRPC service.
func (srv *Server) DeployModule(stream pb.ModuleDeployer_DeployModuleServer) error {
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Printf("Got input from client: %v", in)

	for _, req := range srv.Reqs {
		err = stream.Send(req)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartServer starts the gRPC server goroutine.
func (srv *Server) Start() (*grpc.Server, net.Addr) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Could not listen on ':0'")
	}
	grpcServer := grpc.NewServer()

	pb.RegisterModuleDeployerServer(grpcServer, srv)

	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to start serving gRPC service")
		}
	}()

	return grpcServer, lis.Addr()
}
