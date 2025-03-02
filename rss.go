package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Init client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error with GET request: %w", err)
	}
	req.Header.Set("User-Agent", "aggreGATOR")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error with GET request: %w", err)
	}

	// Handle resp -> data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error with ReadAll: %w", err)
	}
	defer resp.Body.Close()

	// Decode data into xml
	result := RSSFeed{}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error with xml Unmarshal: %w", err)
	}

	// Decode escaped HTML entities
	result.Channel.Description = html.UnescapeString(result.Channel.Description)
	result.Channel.Title = html.UnescapeString(result.Channel.Title)
	for i, item := range result.Channel.Item {
		result.Channel.Item[i].Description = html.UnescapeString(item.Description)
		result.Channel.Item[i].Title = html.UnescapeString(item.Title)
	}

	return &result, nil
}
