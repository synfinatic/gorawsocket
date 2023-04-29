//go:build darwin
// +build darwin

package gorawsocket

import (
	"fmt"
	"net"
	"syscall"
)

func BindDevice(s int, iface *net.Interface) error {
	if err := syscall.SetsockoptInt(s, syscall.IPPROTO_IP, syscall.IP_BOUND_IF, iface.Index); err != nil {
		return fmt.Errorf("unable to IP_BOUND_IF: %s", err.Error())
	}
	return nil
}
