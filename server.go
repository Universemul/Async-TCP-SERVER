package main

import (
	"fmt"
	"net"
)

type TcpServer struct {
	listener net.Listener
	clients  []*Client
}

func (server *TcpServer) Start() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		client := server.addClient(conn)
		go perform_command(client)
	}
}

func (server *TcpServer) Close() {
	server.listener.Close()
}

func (server *TcpServer) Listen(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", PORT)
	listener, err := net.Listen("tcp", addr)
	if err == nil {
		server.listener = listener
	}
	return err
}

func (server *TcpServer) addClient(conn net.Conn) *Client {
	c := &Client{
		conn: conn,
	}
	server.clients = append(server.clients, c)
	return c
}

type Client struct {
	conn net.Conn
}

func (client *Client) Close() {
	client.conn.Close()
}
