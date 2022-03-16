package model

import "net"

type Socket interface {
	Process(net.Conn) error
	Auth(net.Conn) error
	Connect(net.Conn) (net.Conn, error)
	Forward(net.Conn, net.Conn)
}
