package spotify

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethan-a-perry/song-loop/internal/database/data"
)

type Service interface {
	Loop(start, end int)
}

type svc struct {
	userData *data.UserData
}

func NewService(userData *data.UserData) Service {
	return &svc {
		userData: userData,
	}
}

func (s *svc) Loop(start, end int) {
	token, err := s.userData.GetSpotifyToken()
	if err != nil {
		fmt.Print("error: ", err)
		return
	}

	go func() {
		for {
			if err := s.seek(start, token.AccessToken); err != nil {
				fmt.Print("error: ", err)
				return
			}

			duration := end - start
			time.Sleep(time.Duration(duration) * time.Millisecond)
		}
	}()
}

func (s *svc) seek(start int, token string) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", start)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer " + token)

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Request from client failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Spotify returned status: %s", res.Status)
	}

	return nil
}
