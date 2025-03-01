package main

import (
	"context"
	"fmt"
)

// Reset command : for debugging/testing only, delete everything from users table
func handlerReset(s *state, cmd command) error {
	// Ensure that no arguments passed
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset command expect no argument")
	}

	// Call the reset query
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users table: %w", err)
	}
	fmt.Println("Database reset successfully!")

	return nil
}
