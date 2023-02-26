package grpc

import (
	"fmt"
	"net"

	"github.com/tricorder/src/utils/errors"
	"google.golang.org/grpc"
)

// Includes the underlying data structures for serving gRPC RPCs.
type serverFixture struct {
	// A listener to accept client connection.
	// gRPC server needs this to start serving client requests.
	listener net.Listener

	// A server that drives the server process.
	server *grpc.Server

	// The actual address this server listens.
	addr net.Addr
}

// Returns a new serverFixture that listens at the specified port of the localhost.
func newServerFixture(port int) (*serverFixture, error) {
	addrStr := fmt.Sprintf(":%d", prot)
	listener, err := net.Listen("tcp", addrStr)
	if err != nil {
		return nil, errors.Wrap("newing serverFixture", "listen "+addrStr, err)
	}
	return &serverFixture{
		listener: listener,
		server:   grpc.NewServer(),
		addr:     listener.Addr(),
	}, nil
}

// Starts serving gRPC service.
func (f *serverFixture) serve() error {
	return f.server.Serve(f.listener)
}
