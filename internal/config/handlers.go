package config

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/JonahLargen/BlobAggregator/internal/database"
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
	fmt.Printf("Logged in as user %s\n", resp.Name)
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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("scrape feeds command requires a time between requests argument")
	}
	time_between_reqs := cmd.Args[0]
	err := agg(s, time_between_reqs)
	if err != nil {
		return fmt.Errorf("error starting feed aggregation: %w", err)
	}
	return nil
}

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("add feed command requires a feed name and feed URL argument")
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	feed, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking if feed exists: %w", err)
	}
	if feed.ID != uuid.Nil {
		return fmt.Errorf("feed %s already exists", feedURL)
	}
	feed, err = s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
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
	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed %s: %w", feedURL, err)
	}
	fmt.Printf("Added feed %s\n", feed.Url)
	return nil
}

func handlerFeeds(s *State, _ Command) error {
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

func handlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow command requires a url argument")
	}
	feedURL := cmd.Args[0]
	feed, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed %s not found", feedURL)
	}
	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed %s: %w", feedURL, err)
	}
	fmt.Printf("Followed feed %s\n", feedURL)
	return nil
}

func handlerFollowing(s *State, _ Command, user database.User) error {
	following, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching following feeds: %w", err)
	}
	fmt.Println("Following feeds:")
	for _, follow := range following {
		fmt.Printf("- %s\n", follow)
	}
	if len(following) == 0 {
		fmt.Println("You are not following any feeds")
	}
	return nil
}

func handlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("unfollow command requires a url argument")
	}
	feedURL := cmd.Args[0]
	feed, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed %s not found", feedURL)
	}
	_, err = s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed %s: %w", feedURL, err)
	}
	fmt.Printf("Unfollowed feed %s\n", feedURL)
	return nil
}

func handlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	var err error
	if len(cmd.Args) > 0 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit value: %v", err)
		}
	}
	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}
	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2, 2006"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
