package main

import (
	"fmt"
	"os"

	"github.com/ethan-a-perry/song-loop/internal/spotify"
	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
	"github.com/ethan-a-perry/song-loop/internal/store"
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

	store := store.NewStore()
	authService := spotifyauth.NewService(store)
	spotifyService := spotify.NewService(authService)

	api := api {
		config: cfg,
		store: store,
		authService: authService,
		spotifyService: spotifyService,
	}

	router := api.mount()

	if err := api.run(router); err != nil {
		fmt.Printf("Server has failed to start: %s", err)
		os.Exit(1)
	}
}
