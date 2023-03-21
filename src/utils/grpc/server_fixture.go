package grpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/sys"
)

// Includes the underlying data structures for serving gRPC RPCs.
type ServerFixture struct {
	// A listener to accept client connection.
	// gRPC server needs this to start serving client requests.
	listener net.Listener

	// A server that drives the server process.
	Server *grpc.Server

	// The actual address this server listens.
	Addr net.Addr
}

// Returns a new ServerFixture that listens at the specified port of the localhost.
func NewServerFixture(port int) (*ServerFixture, error) {
	listener, addr, err := sys.ListenTCP(port)
	if err != nil {
		return nil, errors.Wrap("newing ServerFixture", "listen to local port", err)
	}
	return &ServerFixture{
		listener: listener,
		Server:   grpc.NewServer(),
		Addr:     addr,
	}, nil
}

// Starts serving gRPC service.
func (f *ServerFixture) Serve() error {
	return f.Server.Serve(f.listener)
}
