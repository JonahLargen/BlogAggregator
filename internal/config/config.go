package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func ReadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	configFilePath := fmt.Sprintf("%s/%s", home, configFileName)
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return &config, nil
}

func (c *Config) SetUser(userName string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	configFilePath := fmt.Sprintf("%s/%s", home, configFileName)
	c.CurrentUserName = userName
	file, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(configFilePath, file, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configFilePath, err)
	}
	return nil
}
