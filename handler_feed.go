package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Ensure that a name and a url was passed in the args
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed command expect two argument -> <name> <url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	// Call function to create feed
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	// Auto create a feed follow
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
	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	// Ensure that no arguments passed
	if len(cmd.Args) != 0 {
		return fmt.Errorf("users command expect no argument")
	}

	// Get feeds info
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	// Handle no feeds
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	// Print results
	fmt.Printf("Found %d feeds:\n", len(feeds))
	fmt.Println("=====================================")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user info: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID: 			%v\n", feed.ID)
	fmt.Printf("* CreatedAt: 		%v\n", feed.CreatedAt)
	fmt.Printf("* UpdatedAt: 		%v\n", feed.UpdatedAt)
	fmt.Printf("* Name: 		%v\n", feed.Name)
	fmt.Printf("* Url: 			%v\n", feed.Url)
	fmt.Printf("* User: 		%v\n", user.Name)
}
