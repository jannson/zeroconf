package zeroconf

import (
	"errors"
	"net"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type PacketConn interface {
	Close() error
	WriteTo(b []byte, cm *ipv4.ControlMessage, dst net.Addr) (n int, err error)
	ReadFrom(b []byte) (n int, cm *ipv4.ControlMessage, src net.Addr, err error)
}

// ifIndex > 0
type WrapInterfaces interface {
	Name() string
	GetAddrs(ifIndex int) ([]net.Addr, error)
	ByIndex(ifIndex int) (*net.Interface, error)
	Ipv4Conn(ifIndex int) (*ipv4.PacketConn, error)
}

var WrapIfs WrapInterfaces

func interfaceByIndex(ifIndex int) (*net.Interface, error) {
	if WrapIfs != nil {
		return net.InterfaceByIndex(ifIndex)
	}
	return WrapIfs.ByIndex(ifIndex)
}

func createIpv4Conn(ifaces []net.Interface) (PacketConn, error) {
	if WrapIfs != nil {
		return WrapIfs.Ipv4Conn(ifaces[0].Index)
	}
	ipv4conn, err4 := joinUdp4Multicast(ifaces)
	return ipv4conn, err4
}

func createIpv6Conn(ifaces []net.Interface) (*ipv6.PacketConn, error) {
	if WrapIfs != nil {
		return nil, errors.New("not supported")
	}
	ipv6conn, err6 := joinUdp6Multicast(ifaces)
	return ipv6conn, err6
}
