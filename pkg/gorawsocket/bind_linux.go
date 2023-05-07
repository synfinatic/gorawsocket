//go:build linux
// +build linux

package gorawsocket

import (
	"fmt"
	"net"
	"syscall"
)

func (rs *RawSocket) Bind(iface *net.Interface) error {
	if err := syscall.SetsockoptString(rs.Fd, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, iface.Name); err != nil {
		return fmt.Errorf("unable to setsockopt(SOL_SOCKET, SO_BINDTODEVICE, %s): %s", iface.Name, err.Error())
	}
	return nil
}
