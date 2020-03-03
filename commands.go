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
	message := ""
	if cmd.name == Quit {
		message = "221 Bye\n"
	} else if cmd.name == Hello {
		message = fmt.Sprintf("250 Pleased to meet you %s\n", cmd.args)
	} else if cmd.name == Welcome {
		message = "220 localhost\n"
	} else if cmd.name == Date {
		dt := time.Now()
		message = dt.Format(DateFormat)
	}
	return []byte(message)
}

func (cmd *Command) Error() []byte {
	message := ""
	if cmd.name == Hello {
		message = fmt.Sprintf("%s Verb need a name\n", cmd.name)
	} else if cmd.name == Date {
		message = "550 Bad state\n"
	} else {
		message = fmt.Sprintf("%s Verb is not recognized\n", cmd.name)
	}
	return []byte(message)
}

func (cmd *Command) isValid(client *Client) bool {
	verb := ""
	if cmd.name == Welcome {
		return true
	}
	for i := range availableVerbs {
		if availableVerbs[i] == cmd.name {
			verb = cmd.name
		}
	}
	if verb == Hello && cmd.args == "" {
		return false
	}
	if verb == "" {
		return false
	}
	if verb == Date {
		if _, ok := client.commands[Hello]; !ok {
			return false
		}
	}
	return true
}

func parse(cmd string) (string, string) {
	tmp := strings.Split(strings.TrimSpace(cmd), " ")
	if len(tmp) > 1 {
		return tmp[0], tmp[1]
	}
	return tmp[0], ""
}
