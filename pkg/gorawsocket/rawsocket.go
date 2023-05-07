package gorawsocket

import (
	"fmt"
	"net"
	"syscall"
)

type RawSocket struct {
	Fd        int
	Interface *net.Interface
	SrcIP     *net.IP
}

func NewRawSocket() (*RawSocket, error) {
	var err error
	rs := &RawSocket{}

	if rs.Fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW); err != nil {
		return rs, fmt.Errorf("unable to open socket(AF_INET, SOCK_RAW, IPPROTO_RAW): %s", err.Error())
	}
	return rs, nil
}

func (rs *RawSocket) BufSize(len int) error {
	// set send buffer size
	if err := syscall.SetsockoptInt(rs.Fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, len); err != nil {
		return fmt.Errorf("unable to setsockopt(SOL_SOCKET, SO_SNDBUF, %d): %s", len, err.Error())
	}
	return nil
}

func (rs *RawSocket) NoRoute(yn bool) error {
	val := 0
	if yn {
		val = 1
	}
	if err := syscall.SetsockoptInt(rs.Fd, syscall.SOL_SOCKET, syscall.SO_DONTROUTE, val); err != nil {
		return fmt.Errorf("unable to setsockopt(SOL_SOCKET, SO_DONTROUTE, %d): %s", val, err.Error())
	}
	return nil
}

func (rs *RawSocket) IncludeIPHeader(yn bool) error {
	val := 0
	if yn {
		val = 1
	}
	if err := syscall.SetsockoptInt(rs.Fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return fmt.Errorf("unable to setsockopt(IPPROTO_IP, IP_HDRINCL, %d): %s", val, err.Error())
	}
	return nil
}

func (rs *RawSocket) Sendmsg4(msg, oob []byte, dstIP net.IP, flags int) (int, error) {
	addr, err := NewSockaddrInet4(dstIP)
	if err != nil {
		return 0, err
	}

	return syscall.SendmsgN(rs.Fd, msg, oob, &addr, flags)
}

func (rs *RawSocket) Sendmsg6(msg, oob []byte, dstIP net.IP, flags int) (int, error) {
	addr, err := NewSockaddrInet6(dstIP)
	if err != nil {
		return 0, err
	}

	return syscall.SendmsgN(rs.Fd, msg, oob, &addr, flags)
}

func (rs *RawSocket) Close() error {
	return syscall.Close(rs.Fd)
}

func NewSockaddrInet4(ip net.IP) (syscall.SockaddrInet4, error) {
	ip4 := ip.To4()
	if nil == ip4 {
		return syscall.SockaddrInet4{}, fmt.Errorf("Not an IPv4 address: %s", ip.String())
	}
	return syscall.SockaddrInet4{
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}, nil
}

func NewSockaddrInet6(ip net.IP) (syscall.SockaddrInet6, error) {
	ip6 := ip.To16()
	if nil == ip6 {
		return syscall.SockaddrInet6{}, fmt.Errorf("Not an IPv6 address: %s", ip.String())
	}
	return syscall.SockaddrInet6{
		Addr: [16]byte{
			ip6[0], ip6[1], ip6[2], ip6[3],
			ip6[4], ip6[5], ip6[6], ip6[7],
			ip6[8], ip6[9], ip6[10], ip6[11],
			ip6[12], ip6[13], ip6[14], ip6[15],
		},
	}, nil
}
