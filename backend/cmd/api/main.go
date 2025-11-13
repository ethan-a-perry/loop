package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/auth"
	"github.com/ethan-a-perry/song-loop/internal/spotify"
)

var token *auth.Token

func loopHandler(w http.ResponseWriter, r *http.Request) {
	// Allow requests from your frontend
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if r.Method == http.MethodOptions {
		// Preflight request, respond 200 OK
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := struct {
		Start int `json:"start"`
		End int `json:"end"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON (request body)", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})

	spotify.Loop(req.Start, req.End, token.AccessToken)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.URL.Query().Get("error")

	if err != "" {
		http.Error(w, "Authorization failed during callback" + err, http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
        return
	}

	auth.EstablishToken(code)

	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, authorized, err := auth.Authenticate()

	if err != nil {
		fmt.Print("Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if authorized {
		token = t
		return
	}

	// Not authorized
	authorizationUrl, err := auth.GetAuthorizationUrl()

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, authorizationUrl, http.StatusFound)
}

func main() {
	if err := auth.InitSpotify(); err != nil {
		return
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/loop", loopHandler)

	fmt.Println("Server running at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}
