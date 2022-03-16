package Src

import (
	"ProxyGo/ProxyLog"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Socket5 struct {
	VER       uint8
	CMD       uint8
	RSV       uint8
	ATYP      uint8
	DSTADDR   []byte
	DSTPORT   uint16
	DSTDOMAIN string
	RAWADDR   *net.TCPAddr
}

// socket5流程  Auth -> Connect -> Forward

func (socket *Socket5) Process(client *net.TCPConn) error {

	return nil
}

// 客户先发言
func (socket *Socket5) Auth(client *net.TCPConn) error {
	buffer := make([]byte, 257)
	// 先读前连个字节的 VER NMETHODS
	n, err := io.ReadFull(client, buffer[:2])
	if err != nil {
		ProxyLog.Err.Println("Socket VER ReadErr ", err)
		return err
	}
	if n != 2 {
		ProxyLog.Err.Println("Socket VER ReadErr ")
		return errors.New("ReadHeadErr")
	}
	ver, nmethods := int(buffer[0]), int(buffer[1])
	// 16 进制
	if ver != 5 {
		ProxyLog.Err.Println("NotSocket5 ")
		return errors.New("NotSocket5")
	}
	// 读支持的方法
	n, err = io.ReadFull(client, buffer[:nmethods])
	if err != nil {
		ProxyLog.Err.Println("ReadMethodErr ")
		return err
	}
	if n != nmethods {
		ProxyLog.Err.Println("ReadMethodErr ")
		return errors.New("ReadMethodNumErr")
	}
	ProxyLog.Info.Println("SocketInfo", " VER ", ver, " NMethods ", nmethods)
	// TODO 服务端选定认证的方法发送给客户端

	// 服务端的响应格式 VER METHOD
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		ProxyLog.Err.Println("SocketServerMethodChooseErr")
		return errors.New("SocketServerMethodChooseErr")
	}
	return nil
}

func (socket *Socket5) Connect(client *net.TCPConn) (net.Conn, error) {
	//	和服务器完成协商过后，告诉服务器具体的目标地址
	//	格式：VER CMD RSV ATYP DST.ADDR DST.PORT
	// VER protocol version: X’05’
	// o CMD
	//		o CONNECT X’01’
	// 		o BIND X’02’
	// 		o UDP ASSOCIATE X’03’
	// o RSV RESERVED
	// o ATYP address type of following address
	// 		o IP V4 address: X’01’
	// 		o DOMAINNAME: X’03’
	// 		o IP V6 address: X’04’
	// o DST.ADDR desired destination address
	//		 ATYP field specifies the type of address
	// o DST.PORT desired destination port in network octet order
	buffer := make([]byte, 128)
	// 先读前4个字节
	n, err := io.ReadFull(client, buffer[:4])
	if err != nil {
		ProxyLog.Err.Println("SocketConnectReadErr ", err)
		return nil, err
	}
	ver, _, atype := int(buffer[0]), int(buffer[1]), int(buffer[3])
	// 协议不支持
	if ver != 5 {
		ProxyLog.Err.Println("SocketVersionErr")
		return nil, errors.New("SocketVersionErr")
	}
	addr := ""
	switch atype {
	// ipv4格式地址
	case 1:
		n, err := io.ReadFull(client, buffer[:4])
		if err != nil {
			ProxyLog.Err.Println("SocketDstAddrReadErr", err)
			return nil, err
		}
		if n != 4 {
			ProxyLog.Err.Println("Invalid ipv4 address")
			return nil, errors.New("invalid ipv4 address")
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buffer[0], buffer[1], buffer[2], buffer[3])
	//	域名格式地址
	case 3:
		// 先读 1 个字节表示长度
		n, err := io.ReadFull(client, buffer[:1])
		if n != 1 || err != nil {
			ProxyLog.Err.Println("SocketDOMAINNAMEReadErr ", err)
			return nil, err
		}
		addLen := int(buffer[0])
		// 再读n个字节的域名
		n, err = io.ReadFull(client, buffer[:addLen])
		if n != addLen || err != nil {
			ProxyLog.Err.Println("SocketDOMAINNAMEReadErr ", err)
			return nil, err
		}
		addr = string(buffer[:addLen])
	//	ipv6格式地址
	case 4:
		//	TODO
		ProxyLog.Info.Println("TODO ipv6")
		return nil, nil
	default:
		ProxyLog.Err.Println("invalid atype")
		return nil, errors.New("invalid atype")
	}
	// 读取 端口 2个字节的无符号数
	n, err = io.ReadFull(client, buffer[:2])
	if n != 2 || err != nil {
		ProxyLog.Err.Println("SocketPortReadErr ", err)
		return nil, err
	}
	// 字节序是大端存储
	port := binary.BigEndian.Uint16(buffer[:2])
	// 得到端口
	dstAddr := fmt.Sprintf("%s:%d", addr, port)
	// TODO 加密，认证
	// 和远端建立连接
	dial, err := net.Dial("tcp", dstAddr)
	if err != nil {
		return nil, err
	}
	// 回复客户端，连接已经就绪
	// 格式：VER REP RSV ATYP BAD.ADDR BAN.PORT
	// TODO 测试用例
	n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dial.Close()
		return nil, errors.New("write rsp: " + err.Error())
	}
	return dial, nil
}

func (socket *Socket5) Forward(client *net.TCPConn) error {
	//TODO implement me
	panic("implement me")
}
