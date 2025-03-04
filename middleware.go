package main

import (
	"context"
	"fmt"

	"github.com/VincNT21/aggreGATOR/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	// Returning a wrapper function
	return func(s *state, cmd command) error {
		// Middleware logic here:
		// Get current user from DB
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't get current user from db: %w", err)
		}

		// pass the user along to the wrapped handler
		return handler(s, cmd, user)
	}
}
