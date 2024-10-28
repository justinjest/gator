package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	json_parser "github.com/justinjest/gator/internal/config"
	"github.com/justinjest/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *json_parser.Config
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
	_, err := json_parser.SetUser(*s.cfg, name)
	if err != nil {
		return err
	}
	fmt.Printf("Username set to %v\n	", name)
	return nil
}

func main() {
	s, err := startUp()
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
	err = c.run(&s, cmd)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}
}

func startUp() (state, error) {
	config, err := json_parser.Read()
	if err != nil {
		return state{}, errors.New("unable to read json")
	}
	var s state
	s.cfg = &config
	db, err := sql.Open("postgres", s.cfg.Db_url)
	if err != nil {
		return state{}, errors.New("unable to open postgres db")
	}
	dbQueries := database.New(db)
	s.db = dbQueries
	return s, nil
}
