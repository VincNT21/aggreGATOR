package main

import (
	"errors"
)

// command struct = each command has a name and zero to many arguments
type command struct {
	Name string
	Args []string
}

// commands struct = to hold all the commands the CLI can handle.
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// commands methods:
// register = to register a new handler function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run = to runs a given command with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return f(s, cmd)
}
