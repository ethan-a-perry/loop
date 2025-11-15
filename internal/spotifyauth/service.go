package spotifyauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethan-a-perry/song-loop/config"
	"github.com/ethan-a-perry/song-loop/internal/database/data"
	"github.com/ethan-a-perry/song-loop/internal/utils"
)

type Service interface {
	Authenticate() (*Token, bool, error)
	EstablishToken(code string) error
	GetAuthorizationUrl() (string, error)
}

type svc struct {
	userData *data.UserData
}

var (
	session *config.Config
	codeVerifier string
)

type Token struct {
	AccessToken string
	TokenType string
	Scope string
	ExpiresAt time.Time
	RefreshToken string
}

func NewService(userData *data.UserData) Service {
	var err error
	session, err = config.LoadConfig()

	if err != nil {
		fmt.Errorf("Failed to load config: %w", err)
	}

	return &svc {
		userData: userData,
	}
}

func (s *svc) GetAuthorizationUrl() (string, error) {
	var err error
	codeVerifier, err = utils.GenerateCodeVerifier(64)
	if err != nil {
		return "", errors.New("Could not generate code verifier")
	}

	codeChallenge := utils.GenerateCodeChallenge(codeVerifier)

	state, err := utils.GenerateCodeVerifier(64)
	if err != nil {
		return "", errors.New("Could not generate state")
	}

	v := url.Values{}
	v.Set("client_id", session.ClientID)
	v.Set("response_type", "code")
	v.Set("redirect_uri", session.RedirectURI)
	v.Set("state", state)
	v.Set("scope", session.Scope)
	v.Set("code_challenge_method", "S256")
	v.Set("code_challenge", codeChallenge)

	return fmt.Sprintf("https://accounts.spotify.com/authorize?%s", v.Encode()), nil
}

func (s *svc) EstablishToken(code string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", session.RedirectURI)
	data.Set("client_id", session.ClientID)
	data.Set("code_verifier", codeVerifier)

	return s.GetToken(data)
}

func (s *svc) RefreshToken(refreshToken string) error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", session.ClientID)

	return s.GetToken(data)
}

func (s *svc) GetToken(data url.Values) error {
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))

	if err != nil {
		return fmt.Errorf("Request creation failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Request from client failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("spotify token exchange failed: %s", body)
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		TokenType string `json:"token_type"`
		Scope string `json:"scope"`
		ExpiresIn int `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	token := Token{
		AccessToken: tr.AccessToken,
		TokenType: tr.TokenType,
		Scope: tr.Scope,
		ExpiresAt: time.Now().Add((time.Duration(tr.ExpiresIn) - 10) * time.Second),
		RefreshToken: tr.RefreshToken,
	}

	if err := s.SaveToken(&token); err != nil {
        return fmt.Errorf("Failed to save token: %w", err)
    }

	return nil
}

func (s *svc) LoadToken() (*Token, error) {
	token, err := s.userData.GetSpotifyToken()

	if err != nil {
		return nil, fmt.Errorf("Failed to get spotify token from database: %w", err)
	}

	return &token, nil
}

func (s *svc) SaveToken(t *Token) error {
	if err := s.userData.UpdateSpotifyToken(&t); err != nil {
        return fmt.Errorf("Failed to save token: %w", err)
    }

    return nil
}

func (s *svc) Authenticate() (*Token, bool, error) {
	t, err := s.LoadToken()

	if err != nil {
		return t, false, nil
	}

	if time.Now().After(t.ExpiresAt) {
		err := s.RefreshToken(t.RefreshToken)

		if err != nil {
			return t, false, fmt.Errorf("Failed to refresh token: %w", err)
		}
	}

	return t, true, nil
}
