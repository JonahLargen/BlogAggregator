package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
		fmt.Printf("User %s does not exist. Please register.\n", userName)
		os.Exit(1)
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
		fmt.Printf("User %s already exists. Please log in instead.\n", name)
		os.Exit(1)
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
