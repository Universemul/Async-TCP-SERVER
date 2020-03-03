package main

import (
	"fmt"
)

// This method will be delete tonight. It will be replaced by a CommandHandler
func perform_command(client *Client) {
	Write(client.conn, "220 localhost\n")
	for {
		cmd, _, err := Read(client.conn)
		if err != nil {
			fmt.Printf("ERROR for client %s : %s\n", client.conn.RemoteAddr(), err)
		}
		if cmd == QUIT {
			err := Write(client.conn, "221 Bye")
			if err != nil {
				fmt.Printf("ERROR for client %s : %s\n", client.conn.RemoteAddr(), err)
			} else {
				break
			}
		}
		client.conn.Write([]byte(cmd))
	}
	client.Close()
}

func main() {
	var server TcpServer
	err := server.Listen(PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	server.Start()
}
