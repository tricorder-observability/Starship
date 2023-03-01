package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/tricorder/src/utils/errors"
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
	const tcp = "tcp"
	addrStr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen(tcp, addrStr)
	if err != nil {
		return nil, errors.Wrap("newing ServerFixture", "listen "+addrStr, err)
	}
	return &ServerFixture{
		listener: listener,
		Server:   grpc.NewServer(),
		Addr:     listener.Addr(),
	}, nil
}

// Starts serving gRPC service.
func (f *ServerFixture) Serve() error {
	return f.Server.Serve(f.listener)
}
