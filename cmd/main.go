package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// ENV file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Configure API
	cfg := config {
		addr: ":8080",
	}

	api := api {
		config: cfg,
	}

	router := api.mount()

	if err := api.run(router); err != nil {
		fmt.Printf("Server has failed to start: %s", err)
		os.Exit(1)
	}
}
