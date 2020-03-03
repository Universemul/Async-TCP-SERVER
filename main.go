package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	port := Port
	if len(args) >= 1 {
		if i, err := strconv.Atoi(args[0]); err == nil {
			port = i
		}
	}
	var server TcpServer
	err := server.Listen(port)
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
