package main

import "ProxyGo/Src"

func main() {
	server := Src.NewServer("Radnom", "127.0.0.1", 8899, 7766)
	server.Start()
}
