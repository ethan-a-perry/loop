package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
	ExpiresAt time.Time `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

type Store struct {
	token SpotifyToken
}

var filepath = "internal/store/token.json"

func NewStore() *Store {
	return &Store{}
}

func(s *Store) Load() (*SpotifyToken, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token SpotifyToken

	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &token, nil
}

func (s *Store) Save(token *SpotifyToken) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	fmt.Println("Saved token successfully.")

	return nil
}
