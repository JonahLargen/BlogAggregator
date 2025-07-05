package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"

	"github.com/JonahLargen/BlobAggregator/internal/database"
	"github.com/JonahLargen/BlobAggregator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	for _, item := range feed.Channel.Item {
		pubDate, err := parsePubDate(item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing pubDate %s: %w", item.PubDate, err)
		}
		_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			// Ignore unique constraint violation on url
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			return fmt.Errorf("error creating post %s: %w", item.Link, err)
		}
		fmt.Printf("Post created: %s (%s)\n", item.Title, item.Link)
	}
	return nil
}

func parsePubDate(pubDate string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC850,   // "Monday, 02-Jan-06 15:04:05 MST"
		time.RFC3339,  // Atom feeds sometimes use this
	}

	// Normalize common timezone abbreviations
	pubDate = strings.Replace(pubDate, "GMT", "UTC", 1)

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, pubDate)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}
