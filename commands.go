package main

import (
	"fmt"
	"time"
)

const (
	Port       = 8001
	Quit       = "QUIT"
	Hello      = "EHLO"
	Date       = "DATE"
	Welcome    = "WELCOME"
	DateFormat = "02-01-2006T15:04:05\n"
)

var availableVerbs = [4]string{Quit, Hello, Date}

type Command interface {
	Display() []byte
	Error() []byte
	IsValid(*Client) bool
	Name() string
}

type QuitCommand struct {
	name, args string
}

type WelcomeCommand struct {
	name, args string
}

type HelloCommand struct {
	name, args string
}

type DateCommand struct {
	name, args string
}

type UnknownCommand struct {
	name, args string
}

func (cmd QuitCommand) Display() []byte {
	return []byte("221 Bye\n")
}

func (cmd WelcomeCommand) Display() []byte {
	return []byte("220 localhost\n")
}

func (cmd DateCommand) Display() []byte {
	dt := time.Now()
	return []byte(dt.Format(DateFormat))
}

func (cmd HelloCommand) Display() []byte {
	return []byte(fmt.Sprintf("250 Pleased to meet you %s\n", cmd.args))
}

func (cmd UnknownCommand) Display() []byte {
	return []byte("")

}

func (cmd QuitCommand) Error() []byte {
	return []byte("\n")
}

func (cmd WelcomeCommand) Error() []byte {
	return []byte("220 localhost\n")
}

func (cmd DateCommand) Error() []byte {
	return []byte("550 Bad state\n")
}

func (cmd HelloCommand) Error() []byte {
	return []byte(fmt.Sprintf("%s Verb need a name\n", cmd.name))
}

func (cmd UnknownCommand) Error() []byte {
	return []byte(fmt.Sprintf("%s Verb is not recognized\n", cmd.name))

}

func (cmd QuitCommand) IsValid(c *Client) bool {
	return true
}

func (cmd HelloCommand) IsValid(c *Client) bool {
	return cmd.args != ""
}

func (cmd DateCommand) IsValid(c *Client) bool {
	if _, ok := c.commands[Hello]; !ok {
		return false
	}
	return true
}

func (cmd UnknownCommand) IsValid(c *Client) bool {
	return false
}

func (cmd WelcomeCommand) IsValid(c *Client) bool {
	return true
}

func (cmd QuitCommand) Name() string {
	return cmd.name
}

func (cmd HelloCommand) Name() string {
	return cmd.name
}

func (cmd DateCommand) Name() string {
	return cmd.name
}

func (cmd UnknownCommand) Name() string {
	return cmd.name
}

func (cmd WelcomeCommand) Name() string {
	return cmd.name
}

func CommandFactoy(verb string, args string) Command {
	switch verb {
	case Quit:
		return QuitCommand{name: verb, args: args}
	case Welcome:
		return WelcomeCommand{name: verb, args: args}
	case Date:
		return DateCommand{name: verb, args: args}
	case Hello:
		return HelloCommand{name: verb, args: args}
	default:
		return UnknownCommand{name: verb, args: args}
	}
}
