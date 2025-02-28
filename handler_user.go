package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login command expect one argument -> username")
	}
	username := cmd.Args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user. Err: %w", err)
	}
	fmt.Printf("User '%v' has been set !\n", username)
	return nil
}
