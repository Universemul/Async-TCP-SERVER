package main

import (
	"fmt"
)

/* TODO:
- Implementer le perform_command pour un client (ajout d'une interface)
	- Implementer un CommandHandler avec dependance entre Command
		Command :
			name
			depends_on
			message
- Faire un vrai Writer avec des methods pour string/bytes
- Faire un vrai Reader


*/

func perform_command(client *Client) {
	Write(client.conn, "220 localhost\n")
	for {
		cmd, _, err := Read(client.conn)
		if err != nil {
			fmt.Printf("ERROR for client %s : %s\n", client.conn.RemoteAddr(), err)
		}
		if cmd == QUIT {
			// Define write method
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
