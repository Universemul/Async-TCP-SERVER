package main

import (
	"fmt"
	"strings"
	"time"
)

type Command struct {
	name string
	args string
}

func (cmd *Command) Display() []byte {
	switch cmd.name {
	case Quit:
		return []byte("221 Bye\n")
	case Hello:
		return []byte(fmt.Sprintf("250 Pleased to meet you %s\n", cmd.args))
	case Welcome:
		return []byte("220 localhost\n")
	case Date:
		dt := time.Now()
		return []byte(dt.Format(DateFormat))
	}
	return []byte("")
}

func (cmd *Command) Error() []byte {
	switch cmd.name {
	case Hello:
		return []byte(fmt.Sprintf("%s Verb need a name\n", cmd.name))
	case Date:
		return []byte("550 Bad state\n")
	default:
		return []byte(fmt.Sprintf("%s Verb is not recognized\n", cmd.name))
	}
}

func (cmd *Command) isValid(client *Client) bool {
	switch cmd.name {
	case "":
		return false
	case Quit, Welcome:
		return true
	case Hello:
		if cmd.args == "" {
			return false
		}
		return true
	case Date:
		if _, ok := client.commands[Hello]; !ok {
			return false
		}
		return true
	default:
		return false
	}
}

func parse(cmd string) (string, string) {
	tmp := strings.Split(strings.TrimSpace(cmd), " ")
	if len(tmp) > 1 {
		return tmp[0], tmp[1]
	}
	return tmp[0], ""
}
