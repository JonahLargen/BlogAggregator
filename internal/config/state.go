package config

import "github.com/JonahLargen/BlogAggregator/internal/database"

type State struct {
	Config *Config
	DB     *database.Queries
}
