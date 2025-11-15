package main

import (
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/spotify"
	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
)

type api struct {
	config config
}

type config struct {
	addr string
}

func (a *api) mount() http.Handler {
	router := http.NewServeMux()

	// Spotify Auth
	spotifyAuthService := spotifyauth.NewService()
	spotifyAuthHandler := spotifyauth.NewHandler(spotifyAuthService)

	router.HandleFunc("/request-spotify", spotifyAuthHandler.RequestSpotify)
	router.HandleFunc("/callback", spotifyAuthHandler.Callback)

	// Spotify
	spotifyService := spotify.NewService()
	spotifyHandler := spotify.NewHandler(spotifyService)

	router.HandleFunc("/loop", spotifyHandler.Loop)

	return router
}

func (a *api) run(router http.Handler) error {
	server := http.Server {
		Addr: a.config.addr,
		Handler: router,
	}

	fmt.Println("Server running at http://127.0.0.1:8080")

	return server.ListenAndServe()
}
