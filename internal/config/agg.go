package config

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/JonahLargen/BlobAggregator/internal/database"
	"github.com/JonahLargen/BlobAggregator/internal/rss"
)

func agg(s *State, time_between_reqs string) error {
	timeBetweenRequests, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("invalid time duration: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("error scraping feeds: %w", err)
		}
	}
}

func scrapeFeeds(s *State) error {
	nextFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching next feed: %w", err)
	}
	fmt.Printf("Scraping feed %s (%s)\n", nextFeed.Name, nextFeed.Url)
	err = s.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}
	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error scraping feed %s: %w", nextFeed.Url, err)
	}
	fmt.Println("Feed scraped successfully")
	maxItems := 5
	for i, item := range feed.Channel.Item {
		if i >= maxItems {
			fmt.Printf("... (%d more items)\n", len(feed.Channel.Item)-maxItems)
			break
		}
		fmt.Printf("- Title: %s\n", item.Title)
	}
	return nil
}
