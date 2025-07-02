package config

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JonahLargen/BlobAggregator/internal/database"
	"github.com/google/uuid"
)

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		if s.Config.CurrentUserName == "" {
			return fmt.Errorf("you must be logged in to run this command")
		}
		user, err := s.DB.GetUserByName(context.Background(), s.Config.CurrentUserName)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error fetching user: %w", err)
		}
		if user.ID == uuid.Nil {
			return fmt.Errorf("user %s does not exist, please register", s.Config.CurrentUserName)
		}
		return handler(s, cmd, user)
	}
}
