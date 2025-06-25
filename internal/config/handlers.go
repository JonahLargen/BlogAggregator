package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JonahLargen/BlobAggregator/internal/database"
	"github.com/JonahLargen/BlobAggregator/internal/rss"
	"github.com/google/uuid"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login command requires a username argument")
	}
	userName := cmd.Args[0]
	resp, err := s.DB.GetUserByName(context.Background(), userName)
	if err == sql.ErrNoRows || resp.ID == uuid.Nil {
		return fmt.Errorf("user %s does not exist, please register", userName)
	}
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	if err := s.Config.SetUser(userName); err != nil {
		return err
	}
	fmt.Printf("Logged in as user %s\n", userName)
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("register command requires a name")
	}
	name := cmd.Args[0]
	resp, err := s.DB.GetUserByName(context.Background(), name)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if resp.ID != uuid.Nil {
		return fmt.Errorf("user %s already exists, please log in instead", name)
	}
	now := time.Now()
	resp, err = s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	if err := s.Config.SetUser(resp.Name); err != nil {
		return fmt.Errorf("error setting user in config: %w", err)
	}
	fmt.Printf("Registered user %s\n", resp.Name)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.DB.ResetAll(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}
	fmt.Printf("Database reset successfully\n")
	return nil
}

func handlerListUsers(s *State, cmd Command) error {
	users, err := s.DB.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error listing users: %w", err)
	}
	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}
	fmt.Println("Users:")
	currentUserName := s.Config.CurrentUserName
	for _, user := range users {
		if user.Name == currentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("agg command requires a feed URL argument")
	}
	feedURL := cmd.Args[0]
	//feedURL := "https://www.wagslane.dev/index.xml"
	fetchFeed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed %s: %w", feedURL, err)
	}
	fmt.Printf("%+v\n", fetchFeed)
	return nil
}

func addFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("add feed command requires a feed name and feed URL argument")
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	user, err := s.DB.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed %s: %w", feedURL, err)
	}
	fmt.Printf("Added feed %s\n", feed.Url)
	return nil
}

func feeds(s *State, _ Command) error {
	feeds, err := s.DB.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error listing feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s | URL: %s | User ID: %s\n",
			feed.Name,
			feed.Url,
			feed.UserName,
		)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found")
	}
	return nil
}
