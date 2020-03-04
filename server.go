package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
)

type Server interface {
	Start() error
	Close()
	Listen() error
}

// Define the current Server
// clients define all the current Connection
// mutex is used to lock/unlock the server when adding/deleting a Client
type TcpServer struct {
	listener net.Listener
	clients  []*Client
	mutex    sync.Mutex
}

// Start the Server and launch a goroutine when a new connection appears
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

// Disconnect a client.  We need to remove this client from the 'active clients'
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

// The core method of the server. It allows the server to read and parse the incomming command and to write a custom message to a specific client
func (server *TcpServer) performCommand(conn net.Conn) {
	client := server.addClient(conn)
	defer server.disconnect(client)
	client.Write(CommandFactoy(Welcome, ""))
	for {
		cmd, err := client.Read()
		if err != nil {
			if cmd == nil || cmd.Name() == "" { // Client has disconnected
				break
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
	commands map[string]int // HashMap of commands asked by the client. The value is the number of time requested
	writer   *bufio.Writer
	reader   *bufio.Reader
}

func (client *Client) Close() {
	client.writer = nil
	client.reader = nil
	client.conn.Close()
	client.conn = nil

}

// Parse the current command. Clean the string and separate the verb and the args (e.g The command EHLO [name])
func (client *Client) Parse(cmd string) (string, string) {
	pattern := regexp.MustCompile(`\s+`)
	cmd = strings.TrimSpace(cmd)
	tmp := strings.Split(pattern.ReplaceAllString(cmd, " "), " ")
	if len(tmp) > 1 {
		return strings.TrimSpace(tmp[0]), strings.TrimSpace(tmp[1])
	}
	return strings.TrimSpace(tmp[0]), ""
}

// Read the command sent by a client
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

// Write something to a client using CommandFactory
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
