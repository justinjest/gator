package main

import (
	"errors"
	"fmt"
	"log"

	json_parser "github.com/justinjest/gator/internal/config"
)

type state struct {
	state *json_parser.Config
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
		return errors.New("no arguments passed to login function")
	}
	if len(cmd.args) >= 2 {
		return errors.New("login takes exactly one argument")
	}
	name := cmd.args[0]
	_, err := json_parser.SetUser(*s.state, name)
	if err != nil {
		return err
	}
	fmt.Printf("Username set to %v", name)
	return nil
}

func main() {
	var s state
	config, err := json_parser.Read()
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	s.state = &config

}
