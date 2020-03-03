package main

import (
	"fmt"
)

func main() {
	var server TcpServer
	err := server.Listen(Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = server.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
