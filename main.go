package main

import (
	"log"
	"os"

	"github.com/VincNT21/aggreGATOR/internal/config"
)

// state struct = to give handlers access to the application state (cfg, db connection...)
type state struct {
	cfg *config.Config
}

func main() {
	// Read the config file
	configData, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Store the config in a new instance of the state struct
	programState := &state{
		cfg: &configData,
	}

	// Create a new instance of the commands structs
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register a handler function for the login command
	cmds.register("login", handlerLogin)

	// Get the command-line arguments passed in by the user
	args := os.Args
	// If fewer than 2 arguments => error (first is auto the program name, second is a command name)
	if len(args) < 2 {
		log.Fatalf("no command / not enough arguments specified")
		return
	}
	// Create a command instance based on command/args given
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	// Run the command given
	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
