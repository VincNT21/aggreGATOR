package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

// Register command : Create a new user in db and set cfg properly
func handlerRegister(s *state, cmd command) error {
	// Ensure that a name was passed in the args
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register command expect one argument -> username")
	}
	username := cmd.Args[0]

	// Create a new user in the database
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	// Set the current user in the config to the given name
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	// Print validation message + debug user info
	fmt.Printf("User '%v' created\n", username)
	printUser(user)

	return nil
}

// Login command : Check if user in db and set config properly
func handlerLogin(s *state, cmd command) error {
	// Ensure that a name was passed in the args
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login command expect one argument -> username")
	}
	username := cmd.Args[0]

	// Check if user is in the database
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user doesn't exist in database: %w", err)
	}

	// Set the current user in the config to the given name
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user. Err: %w", err)
	}

	// Print validation message
	fmt.Printf("User '%v' has been set !\n", username)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
