package main

import (
	"errors"
	"fmt"
	"os"

	json_parser "github.com/justinjest/gator/internal/config"
)

type state struct {
	config *json_parser.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	method map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.method[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	err := c.method[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("user name requried for login function")
	}
	if len(cmd.args) >= 2 {
		return errors.New("login takes exactly one argument")
	}
	name := cmd.args[0]
	_, err := json_parser.SetUser(*s.config, name)
	if err != nil {
		return err
	}
	fmt.Printf("Username set to %v\n", name)
	return nil
}

func main() {
	config, err := json_parser.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	c := commands{
		method: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	if len(os.Args) < 2 {
		err = errors.New("too few cmdline arguments")
	}
	if err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
	cmd := command{
		os.Args[1],
		os.Args[2:],
	}
	var s state
	s.config = &config
	err = c.run(&s, cmd)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}
}
