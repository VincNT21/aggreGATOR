package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	// Ensure that a name and a url was passed in the args
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed command expect two argument -> <name> <url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	// Get current user from DB
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get current user from db: %w", err)
	}

	// Call function
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	// Print results
	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID: 			%v\n", feed.ID)
	fmt.Printf("* CreatedAt: 	%v\n", feed.CreatedAt)
	fmt.Printf("* UpdatedAt: 	%v\n", feed.UpdatedAt)
	fmt.Printf("* Name: 		%v\n", feed.Name)
	fmt.Printf("* Url: 			%v\n", feed.Url)
	fmt.Printf("* UserID: 		%v\n", feed.UserID)
}
