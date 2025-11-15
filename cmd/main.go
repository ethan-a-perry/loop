package main

import (
	"fmt"
	"os"
)

func main() {
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
