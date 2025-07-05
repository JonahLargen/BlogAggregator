package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JonahLargen/BlogAggregator/internal/config"
	"github.com/JonahLargen/BlogAggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	dbQueries := database.New(db)

	state := config.State{
		Config: cfg,
		DB:     dbQueries,
	}

	commands := config.NewCommands()
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Not enough arguments provided. Usage: <command> [args]")
		os.Exit(1)
	}

	commandName := args[1]
	commandArgs := args[2:]

	err = commands.Run(&state, config.Command{
		Name: commandName,
		Args: commandArgs,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command '%s': %v\n", commandName, err)
		os.Exit(1)
	}
}
