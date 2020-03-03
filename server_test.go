package main

import "testing"

func testvalid(t *testing.T, c *Client, name string) bool {
	cmd := Command{name: name}
	valid := cmd.isValid(c)
	return valid
}

func TestWelcome(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	valid := testvalid(t, c, Welcome)
	if !valid {
		t.Errorf("Valid Welcom = %t; want true", valid)
	}
}

func TestQuit(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	valid := testvalid(t, c, Quit)
	if !valid {
		t.Errorf("Valid Quit = %t; want true", valid)
	}
}

func TestDateWihoutHello(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	valid := testvalid(t, c, Date)
	if valid {
		t.Errorf("Valid Date Without Hello = %t; want false", valid)
	}
}

func TestDateWithHello(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	c.commands[Hello] += 1
	valid := testvalid(t, c, Date)
	if !valid {
		t.Errorf("Valid Date With Hello = %t; want true", valid)
	}
}

func TestHelloWithoutArgs(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	cmd := Command{name: Hello}
	valid := cmd.isValid(c)
	if valid {
		t.Errorf("Valid Hello = %t; want false", valid)
	}
}

func TestHelloWithArgs(t *testing.T) {
	c := &Client{
		commands: make(map[string]int),
	}
	cmd := Command{name: Hello, args: "David"}
	valid := cmd.isValid(c)
	if !valid {
		t.Errorf("Valid Hello = %t; want true", valid)
	}
}
