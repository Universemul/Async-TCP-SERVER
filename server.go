package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
)

type TcpServer struct {
	listener net.Listener
	clients  []*Client
	mutex    sync.Mutex
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
	if err != nil {
		return err
	}
	server.listener = listener
	return nil
}

func (server *TcpServer) addClient(conn net.Conn) *Client {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	c := &Client{
		conn:     conn,
		commands: make(map[string]int),
		writer:   bufio.NewWriter(conn),
		reader:   bufio.NewReader(conn),
	}
	server.clients = append(server.clients, c)
	return c
}

func (server *TcpServer) disconnect(client *Client) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
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
	client.Write(CommandFactoy(Welcome, ""))
	for {
		cmd, err := client.Read()
		if err != nil {
			if cmd.Name() == "" {
				break // Client has disconnected
			} else {
				fmt.Printf("%s", cmd.Error())
				client.Write(cmd)
			}
		} else {
			if cmd.Name() == Quit {
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
	writer   *bufio.Writer
	reader   *bufio.Reader
}

func (client *Client) Close() {
	client.writer = nil
	client.reader = nil
	client.conn.Close()
	client.conn = nil

}

func (client *Client) Parse(cmd string) (string, string) {
	tmp := strings.Split(strings.TrimSpace(cmd), " ")
	if len(tmp) > 1 {
		return tmp[0], tmp[1]
	}
	return tmp[0], ""
}

func (client *Client) Read() (Command, error) {
	var buffer bytes.Buffer
	for {
		line, isPrefix, err := client.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		buffer.Write(line)
		if !isPrefix {
			break
		}
	}
	verb, args := client.Parse(buffer.String())
	return CommandFactoy(verb, args), nil
}

func (client *Client) Write(cmd Command) error {
	f := cmd.Display
	if !cmd.IsValid(client) {
		f = cmd.Error
	} else {
		client.commands[cmd.Name()] += 1
	}
	_, err := client.writer.Write(f())
	if err == nil {
		client.writer.Flush()
	}
	return err
}
