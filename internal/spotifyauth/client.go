package spotifyauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethan-a-perry/song-loop/internal/store"
)

func GetAccessToken(data url.Values) (*store.SpotifyToken, error) {
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request from client failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify token exchange failed")
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		TokenType string `json:"token_type"`
		Scope string `json:"scope"`
		ExpiresIn int `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &store.SpotifyToken{
		AccessToken: tr.AccessToken,
		TokenType: tr.TokenType,
		Scope: tr.Scope,
		ExpiresAt: time.Now().Add((time.Duration(tr.ExpiresIn) - 10) * time.Second),
		RefreshToken: tr.RefreshToken,
	}, nil
}
