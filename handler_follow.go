package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	// Ensure that an url was passed in the args
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow command expect one argument -> url")
	}
	url := cmd.Args[0]

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
	fmt.Println("New feed follow created:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	// Ensure that no arguments passed
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset command expect no argument")
	}

	// Get results from query
	feedsList, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feeds follow list: %w", err)
	}

	if len(feedsList) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	// Print results
	fmt.Println("============")
	fmt.Printf("Current user '%s' follow %d feeds:\n", user.Name, len(feedsList))
	for _, feed := range feedsList {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	fmt.Println()

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	// Ensure that an url was passed in the args
	if len(cmd.Args) != 1 {
		return fmt.Errorf("unfollow command expect one argument -> url")
	}
	url := cmd.Args[0]

	// Get feed by url
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed by given url: %w", err)
	}

	// Call query function
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
