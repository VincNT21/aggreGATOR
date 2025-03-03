package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	// Ensure that an url was passed in the args
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow command expect one argument -> url")
	}
	url := cmd.Args[0]

	// Get current user info
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get current user from db: %w", err)
	}

	// Get feed info by given URL
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed by url from db: %w", err)
	}

	// Create feed follow
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	// Print results
	fmt.Println("New feed follow set!")
	fmt.Printf("Feed name: %s\n", feedFollow.FeedName)
	fmt.Printf("User name: %s\n", feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	// Ensure that no arguments passed
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset command expect no argument")
	}

	// Get current User id
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user info: %w", err)
	}
	userId := user.ID

	// Get results from query
	feedsList, err := s.db.GetFeedFollowForUser(context.Background(), userId)
	if err != nil {
		return fmt.Errorf("couldn't get feeds follow list: %w", err)
	}

	// Print results
	fmt.Println("============")
	fmt.Printf("Current user '%s' follow %d feeds\n", user.Name, len(feedsList))
	for _, feed := range feedsList {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	fmt.Println()

	return nil
}
