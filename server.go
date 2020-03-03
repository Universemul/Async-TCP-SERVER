package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type TcpServer struct {
	listener net.Listener
	clients  []*Client
}

func (server *TcpServer) Start() error {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			return err
		}
		go server.performCommand(conn)
	}
}

func (server *TcpServer) Close() {
	server.listener.Close()
}

func (server *TcpServer) Listen(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err == nil {
		server.listener = listener
	}
	return err
}

func (server *TcpServer) addClient(conn net.Conn) *Client {
	c := &Client{
		conn:     conn,
		commands: make(map[string]int),
	}
	server.clients = append(server.clients, c)
	return c
}

func (server *TcpServer) disconnect(client *Client) {
	for i, check := range server.clients {
		if check == client {
			server.clients = append(server.clients[:i], server.clients[i+1:]...)
		}
	}
	client.Close()
}

func (server *TcpServer) performCommand(conn net.Conn) {
	client := server.addClient(conn)
	defer server.disconnect(client)
	client.Write(Command{name: Welcome})
	for {
		cmd, err := client.Read()
		if err != nil {
			if cmd.name == "" {
				break // Client has disconnected
			} else {
				fmt.Printf("%s", cmd.Error())
				client.Write(cmd)
			}
		} else {
			if cmd.name == Quit {
				err := client.Write(cmd)
				if err != nil {
					fmt.Printf("%s", err)
					break
				} else {
					break
				}
			} else {
				client.Write(cmd)
			}
		}
	}
}

type Client struct {
	conn     net.Conn
	commands map[string]int
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Read() (Command, error) {
	reader := bufio.NewReader(client.conn)
	var buffer bytes.Buffer
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return Command{}, err
		}
		buffer.Write(line)
		if !isPrefix {
			break
		}
	}
	verb, attrs := parse(buffer.String())
	return Command{name: verb, args: attrs}, nil
}

func (client *Client) Write(cmd Command) error {
	writer := bufio.NewWriter(client.conn)
	f := cmd.Display
	if !cmd.isValid(client) {
		f = cmd.Error
	} else {
		client.commands[cmd.name] += 1
	}
	_, err := writer.Write(f())
	if err == nil {
		writer.Flush()
	}
	return err
}
