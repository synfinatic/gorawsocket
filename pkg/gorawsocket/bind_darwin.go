//go:build darwin
// +build darwin

package gorawsocket

import (
	"fmt"
	"net"
	"syscall"
)

func (rs *RawSocket) Bind(iface *net.Interface) error {
	if err := syscall.SetsockoptInt(rs.Fd, syscall.IPPROTO_IP, syscall.IP_BOUND_IF, iface.Index); err != nil {
		return fmt.Errorf("unable to setsockopt(IPPROTO_IP, IP_BOUND_IF, %d): %s", iface.Index, err.Error())
	}
	return nil
}
