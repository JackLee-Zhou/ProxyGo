package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8899")
	if err != nil {
		return
	}
	accept, err := listen.Accept()
	if err != nil {
		return
	}
	data := []byte{0011}
	fmt.Println("Write ", data)
	accept.Write(data)

}
