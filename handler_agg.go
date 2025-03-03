package main

import (
	"context"
	"fmt"
)

// Agg command
func handlerAgg(s *state, cmd command) error {
	// Ensure that no arguments passed
	if len(cmd.Args) != 0 {
		return fmt.Errorf("users command expect no argument")
	}

	// Get RSSfeed
	testURL := "https://feedfry.com/rss/11eff7941c1af1ddbf20d1719fcf69e6"
	rssfeed, err := fetchFeed(context.Background(), testURL)
	if err != nil {
		return fmt.Errorf("couldn't get rssfeed from %v: %w", testURL, err)
	}

	fmt.Println("----RSS FEED----")
	fmt.Printf("Channel title: %v\n", rssfeed.Channel.Title)
	fmt.Printf("Channel link: %v\n", rssfeed.Channel.Link)
	fmt.Printf("Channel description: %v\n", rssfeed.Channel.Description)
	for _, item := range rssfeed.Channel.Item {
		fmt.Println("----RSS ITEM----")
		fmt.Printf("Item title: %v\n", item.Title)
		fmt.Printf("Item link: %v\n", item.Link)
		fmt.Printf("Item description: %v\n", item.Description)
		fmt.Printf("Item PubDate: %v\n", item.PubDate)
	}

	return nil
}
