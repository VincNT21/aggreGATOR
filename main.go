package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/VincNT21/aggreGATOR/internal/config"
	"github.com/VincNT21/aggreGATOR/internal/database"

	// Import a SQL driver but not use it in code (told by the _ in front)
	_ "github.com/lib/pq"
)

// state struct = to give handlers access to the application state (cfg, db connection...)
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Read the config file
	configData, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", configData.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	// Create a new *database.queries
	dbQueries := database.New(db)

	// Store the config and dbQueries in a new instance of the state struct
	programState := &state{
		db: dbQueries,
		cfg: &configData
	}

	// Create a new instance of the commands structs
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register handler functions : login and register command
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

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
