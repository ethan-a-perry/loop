package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PlaybackState struct {
	Item struct {
		ID string `json:"id"`
	} `json:"item"`
	IsPlaying bool `json:"is_playing"`
	Device struct {
		IsActive bool `json:"is_active"`
	} `json:"device"`
}

func GetPlaybackState(accessToken string) (*PlaybackState, error) {
	url := "https://api.spotify.com/v1/me/player"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + accessToken)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request from client failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify returned status: %s", res.Status)
	}

	var playbackState PlaybackState
	if err := json.NewDecoder(res.Body).Decode(&playbackState); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response body: %w", err)
	}

	fmt.Println(playbackState)

	return &playbackState, nil
}

func Seek(start int, accessToken string) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", start)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + accessToken)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request from client failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("spotify returned status: %s", res.Status)
	}

	return nil
}
