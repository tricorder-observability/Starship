package sys

import (
	"fmt"
	"net"

	"github.com/tricorder/src/utils/errors"
)

// Consts for protocol types.
const (
	LOCALHOST = "localhost"
	TCP       = "tcp"
	HTTP      = "http"
)

// PortAddr returns a string as the address of localhost with port.
func PortAddr(port int) string {
	return fmt.Sprintf(":%d", port)
}

// HostPortAddr returns a string as host:port.
func HostPortAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Returns a listener and its address at the specified port of the localhost.
func ListenTCP(port int) (net.Listener, net.Addr, error) {
	addrStr := PortAddr(port)
	listener, err := net.Listen(TCP, addrStr)
	if err != nil {
		return nil, nil, errors.Wrap("newing ServerFixture", "listen "+addrStr, err)
	}
	return listener, listener.Addr(), nil
}
