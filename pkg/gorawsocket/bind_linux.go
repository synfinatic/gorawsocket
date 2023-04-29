//go:build linux
// +build linux

package gorawsocket

import (
	"fmt"
	"net"
	"syscall"
)

func BindDevice(s int, iface *net.Interface) error {
	if err := syscall.SetsockoptString(s, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, iface.Name); err != nil {
		return fmt.Errorf("unable to SO_BINDTODEVICE: %s", err.Error())
	}
	return nil
}
