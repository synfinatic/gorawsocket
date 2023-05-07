//go:build freebsd || dragonfly || netbsd || openbsd
// +build freebsd dragonfly netbsd openbsd

package gorawsocket

import (
	"fmt"
	"net"
	"syscall"
)

func (rs *RawSocket) Bind(iface *net.Interface) error {
	addrs, err := iface.Addrs()
	if err != nil {
		return fmt.Errorf("unable to get addresses: %s", err.Error())
	}
	var addr net.IP
	for _, a := range addrs {
		if addr, _, err = net.ParseCIDR(a.String()); err != nil {
			return fmt.Errorf("unable to parse %s: %s", a.String(), err.Error()
		}
		if addr.To4() != nil {
			break
		}
	}
	if addr == nil {
		return fmt.Errorf("unable to bind to %s", iface.Name)
	}

	sa := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{addr[0], addr[1], addr[2], addr[3]},
	}
	if err = syscall.Bind(rs.Fd, &sa); err != nil {
		return fmt.Errorf("unable to bind: %s", err.Error())
	}
	return nil
}
