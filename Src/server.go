package Src

import (
	"ProxyGo/ProxyLog"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	// 加密方式
	CryptoType string
	//	远端VPS地址
	VpsAddress string
	VpsPort    int
	//	本地监听端口
	LocalPort int
	ProxyType string
}

func NewServer(Crypto string, VpsAddress string, VpsPort int, LocalPort int, ProxyType string) *Server {
	// 实例化 Server 对象
	return &Server{
		CryptoType: Crypto,
		VpsAddress: VpsAddress,
		VpsPort:    VpsPort,
		LocalPort:  LocalPort,
		ProxyType:  ProxyType,
	}
}
func handlerProxy(client *net.TCPConn, s *Server) {
	if client == nil {
		return
	}
	// TODO 对信息先加密之后在连接远端服务器
	// 连接真正的远端的Vps服务器，并设置超时
	dial, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.VpsAddress, s.VpsPort), time.Minute)
	if err != nil {
		ProxyLog.Err.Println("DialVpsErr ", err)
		return
	} else {
		// 排除 有err的情况，defer不应该出现
		defer dial.Close()
	}
	ProxyLog.Info.Println("[Dial] ", dial.LocalAddr().String(), " ", dial.RemoteAddr().String())

	// 同步
	w := sync.WaitGroup{}
	w.Add(2)
	// 本地复制到远端
	go func() {
		defer w.Done()
		Copy(client, dial)
	}()
	//	远端复制到本地
	go func() {
		defer w.Done()
		Copy(dial, client)
	}()
	w.Wait()
}
func (s *Server) Start() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.LocalPort))
	if err != nil {
		ProxyLog.Err.Println("ResolveTCPAddr ", err)
		return
	}
	// 监听本地指定端口
	listener, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	if err != nil {
		ProxyLog.Err.Println("listenLocalErr ", err)
		return
	}
	for {
		// 先连接客户端，客户端连接成功后在连接远端Vps
		accept, err := listener.AcceptTCP()
		if err != nil {
			ProxyLog.Err.Println("AcceptErr ", err)
			return
		}
		ProxyLog.Info.Println("[ServerAccept]: ", " Local: ", accept.LocalAddr().String(),
			" Remote: ", accept.RemoteAddr().String(),
		)
		time.Sleep(time.Second * 10)
		//TODO HTTP socket5代理

		// 获取客户端的代理请求

		// 单独开启协程处理啊远端Vps的连接
		go handlerProxy(accept, s)
	}

}
func (s *Server) EndServer() {
	//TODO implement me
	panic("implement me")
}
