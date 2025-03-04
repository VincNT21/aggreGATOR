package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/VincNT21/aggreGATOR/internal/database"
	"github.com/google/uuid"
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

func scrapeFeeds(s *state) {
	// Get the next feed to fetch from DB
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("couldn't get next feed to fetch: %w", err)
	}

	// Call scrapeFeed on next feed
	scrapeFeed(s.db, nextFeed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	// Mark the next feed as fetched
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fetched: %w", feed.Name, err)
	}

	// Fetch the feed using the URL
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldn't collect feed: %w", err)
	}

	// Store results to posts table
	fmt.Println()
	fmt.Printf("=========RSS Fetch Result from %s========\n", rssFeed.Channel.Title)

	for _, item := range rssFeed.Channel.Item {
		// parse pubdate
		pubDate, err := parseRSSDate(item.PubDate)
		if err != nil {
			log.Println(err)
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("couldn't create post: %v", err)
			continue
		}
		fmt.Printf(" * %s added to posts table\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))

}

func parseRSSDate(dateStr string) (time.Time, error) {
	// List of possible formats to try
	formats := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
		"2006-01-02T15:04:05-07:00",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("couldn't parse date: %s", dateStr)
}
