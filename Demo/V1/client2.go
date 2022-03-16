package main

import (
	"fmt"
	"net"
)

func main() {
	dial, err := net.Dial("tcp", ":7766")
	if err != nil {
		return
	}
	data := []byte{}
	dial.Read(data)
	fmt.Println(data)
}
