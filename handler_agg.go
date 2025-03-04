package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
)

// Agg command
func handlerAgg(s *state, cmd command) error {
	// Ensure that one arguments is passed : duration
	if len(cmd.Args) != 1 {
		return fmt.Errorf("agg command expect one argument -> duration")
	}

	// Parse duration
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse duration time: %w", err)
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetweenRequests)

	// Init a ticker
	ticker := time.NewTicker(timeBetweenRequests)

	// Loop for fetching RSSfeeds
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	/*
		OLD PRINTING FEED PART

			testURL := "https://www.sudeducation.org/feed/"
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
	*/

}

func scrapeFeeds(s *state) error {
	// Get the next feed to fetch from DB
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	// Call scrapeFeed on next feed
	err = scrapeFeed(s.db, nextFeed)
	if err != nil {
		return err
	}
	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) error {
	// Mark the next feed as fetched
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed %s fetched: %w", feed.Name, err)
	}

	// Fetch the feed using the URL
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't collect feed: %w", err)
	}

	// Print the results
	fmt.Println()
	fmt.Printf("=========RSS Fetch Result from %s========\n", rssFeed.Channel.Title)
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf(" * %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

	return nil
}
