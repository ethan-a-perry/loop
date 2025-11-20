package spotifyauth

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ethan-a-perry/song-loop/internal/store"
	"github.com/ethan-a-perry/song-loop/internal/utils"
)

type Service struct {
	store *store.Store
}

var codeVerifier string

func NewService(store *store.Store) *Service {
	return &Service {
		store: store,
	}
}

func (s *Service) GetValidToken() (*store.SpotifyToken, error) {
	token, err := s.store.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	if token == nil {
		return nil, fmt.Errorf("failed to retrieve token. must authenticate with Spotify first.")
	}

	if time.Now().After(token.ExpiresAt) {
		token, err = s.RefreshToken(token.RefreshToken)
		if err != nil {
			return &store.SpotifyToken{}, fmt.Errorf("failed to refresh ticket: %w", err)
		}

		if err := s.store.Save(token); err != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", err)
		}
	}

	return token, nil
}

func (s *Service) GetAuthorizationUrl() (string, error) {
	var err error
	codeVerifier, err = utils.GenerateCodeVerifier(64)
	if err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	state, err := utils.GenerateCodeVerifier(64)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	codeChallenge := utils.GenerateCodeChallenge(codeVerifier)

	v := url.Values{}
	v.Set("client_id", os.Getenv("CLIENT_ID"))
	v.Set("response_type", "code")
	v.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	v.Set("state", state)
	v.Set("scope", os.Getenv("SCOPE"))
	v.Set("code_challenge_method", "S256")
	v.Set("code_challenge", codeChallenge)

	return fmt.Sprintf("https://accounts.spotify.com/authorize?%s", v.Encode()), nil
}

func (s *Service) EstablishToken(code string) error {
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	v.Set("client_id", os.Getenv("CLIENT_ID"))
	v.Set("code_verifier", codeVerifier)

	token, err := GetAccessToken(v)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	if err := s.store.Save(token); err != nil {
		return fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return nil
}

func (s *Service) RefreshToken(refreshToken string) (*store.SpotifyToken, error) {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)
	v.Set("client_id", os.Getenv("CLIENT_ID"))

	token, err := GetAccessToken(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return token, nil
}
