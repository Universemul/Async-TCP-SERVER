package main

import (
	"bytes"
	"testing"
)

func testvalid(c *Client, verb string, args string) bool {
	cmd := CommandFactoy(verb, args)
	valid := cmd.IsValid(c)
	return valid
}

func TestWelcome(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	valid := testvalid(c, Welcome, "")
	if !valid {
		t.Errorf("Valid Welcom = %t; want true", valid)
	}
}

func TestQuit(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	valid := testvalid(c, Quit, "")
	if !valid {
		t.Errorf("Valid Quit = %t; want true", valid)
	}
}

func TestDateWihoutHello(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	valid := testvalid(c, Date, "")
	if valid {
		t.Errorf("Valid Date Without Hello = %t; want false", valid)
	}
}

func TestDateWithHello(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	c.commands[Hello] += 1
	valid := testvalid(c, Date, "")
	if !valid {
		t.Errorf("Valid Date With Hello = %t; want true", valid)
	}
}

func TestHelloWithoutArgs(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	cmd := CommandFactoy(Hello, "")
	valid := cmd.IsValid(c)
	if valid {
		t.Errorf("Valid Hello = %t; want false", valid)
	}
}

func TestHelloWithArgs(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	cmd := CommandFactoy(Hello, "David")
	valid := cmd.IsValid(c)
	if !valid {
		t.Errorf("Valid Hello = %t; want true", valid)
	}
}

func TestQuitDisplay(t *testing.T) {
	valid_msg := []byte("221 Bye\n")
	cmd := CommandFactoy(Quit, "")
	fmtCmd := cmd.Display()
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Quit Display = %s; want %s", fmtCmd, valid_msg)
	}
}

func TestQuitError(t *testing.T) {
	valid_msg := []byte("\n")
	cmd := CommandFactoy(Quit, "")
	fmtCmd := cmd.Error()
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Quit Error = %s; want %s", fmtCmd, valid_msg)
	}
}

func TestWelcomeDisplay(t *testing.T) {
	valid_msg := []byte("220 localhost\n")
	cmd := CommandFactoy(Welcome, "")
	fmtCmd := cmd.Display()
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Welcome Display = %s; want %s", fmtCmd, valid_msg)
	}
}

func TestWelcomeError(t *testing.T) {
	valid_msg := []byte("\n")
	cmd := CommandFactoy(Welcome, "")
	fmtCmd := cmd.Error()
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Welcome Error = %s; want %s", fmtCmd, valid_msg)
	}
}

func TestHelloDisplay(t *testing.T) {
	valid_msg := []byte("250 Pleased to meet you David\n")
	cmd := CommandFactoy(Hello, "David")
	fmtCmd := cmd.Display()
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Hello Display = %s; want %s", fmtCmd, valid_msg)
	}
}

func TestHelloError(t *testing.T) {
	valid_msg := []byte("EHLO Verb need a name\n")
	cmd := CommandFactoy(Hello, "")
	fmtCmd := cmd.Error()
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Hello Error = %s; want %s", fmtCmd, valid_msg)
	}
}

func TestDateDisplay(t *testing.T) {
	valid_msg := []byte("550 Bad state\n")
	cmd := CommandFactoy(Date, "")
	fmtCmd := cmd.Display()
	// We check that Date command does not return "550 Bad state\n". We assume here the command is valid (The client sent the commande EHLO [name] before)
	if bytes.Compare(fmtCmd, valid_msg) == 0 {
		t.Errorf("Message Date Display = %s; want a formatted date like %s", fmtCmd, DateFormat)
	}
}

func TestDateError(t *testing.T) {
	valid_msg := []byte("550 Bad state\n")
	cmd := CommandFactoy(Date, "")
	fmtCmd := cmd.Error()
	// We check that Date command returns "550 Bad state\n". We assume here the command is not valid (The client does not send the commande EHLO [name] before)
	if bytes.Compare(fmtCmd, valid_msg) > 0 {
		t.Errorf("Message Date Error = %s; want %s", fmtCmd, valid_msg)
	}
}
