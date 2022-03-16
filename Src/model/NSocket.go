package model

import "net"

type Socket interface {
	Process(*net.TCPConn) error
	Auth(*net.TCPConn) error
	Connect(*net.TCPConn) (net.Conn, error)
	Forward(*net.TCPConn) error
}
