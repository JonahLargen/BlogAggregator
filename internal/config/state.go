package config

import "github.com/JonahLargen/BlobAggregator/internal/database"

type State struct {
	Config *Config
	DB     *database.Queries
}
