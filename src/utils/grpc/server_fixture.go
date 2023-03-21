package grpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/sys"
)

// Includes the underlying data structures for serving gRPC RPCs.
type ServerFixture struct {
	// A listener listening on a local port for accepting client connections.
	Listener *sys.TCPServer

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
		Listener: listener,
		Server:   grpc.NewServer(),
		Addr:     addr,
	}, nil
}

// Serve starts serving gRPC service. This is a blocking operation.
func (f *ServerFixture) Serve() error {
	return f.Server.Serve(f.Listener)
}
