package main

import (
	"bufio"
	"bytes"
	"net"
	"strings"
)

func parse(cmd string) (string, string) {
	trimCmd := strings.TrimSpace(cmd)
	return trimCmd, ""
}

func Read(conn net.Conn) (string, string, error) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		buffer.Write(line)
		if !isPrefix {
			break
		}
	}
	cmd, attr := parse(buffer.String())
	return cmd, attr, nil
}

func Write(conn net.Conn, cmd string) error {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(cmd)
	return err
}
