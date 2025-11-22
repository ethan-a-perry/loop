package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/spotify"
	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
	"github.com/ethan-a-perry/song-loop/internal/store"
)

type api struct {
	config config
}

type config struct {
	addr string
}

func (a *api) mount() http.Handler {
	router := http.NewServeMux()

	store := store.NewStore()

	// Auth
	authService := spotifyauth.NewService(store)
	authHandler := spotifyauth.NewHandler(*authService)

	router.HandleFunc("/api/spotify/connect", authHandler.Connect)
	router.HandleFunc("/api/spotify/callback", authHandler.Callback)

	// Loop
	spotifyService := spotify.NewService(authService)
	spotifyHandler := spotify.NewHandler(*spotifyService)

	router.HandleFunc("POST /api/spotify/loop", spotifyHandler.Loop)
	router.HandleFunc("/api/spotify/loop/stop", spotifyHandler.StopLoop)

	// App
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

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
