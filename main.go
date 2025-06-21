package main

import (
	"encoding/json"
	"fmt"

	"github.com/JonahLargen/BlobAggregator/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	cfg.SetUser("JonahLargen")
	cfg, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("Error reading config after setting user: %v\n", err)
		return
	}
	jsonBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling config:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
