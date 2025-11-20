package spotify

import (
	"fmt"
	"time"

	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
)

type Service struct {
	auth *spotifyauth.Service
}

func NewService(auth *spotifyauth.Service) *Service {
	return &Service {
		auth: auth,
	}
}

func (s *Service) Loop(start, end int) error {
	token, err := s.auth.GetValidToken()
	if err != nil || token == nil {
		return fmt.Errorf("failed to get valid token")
	}

	go func() {
		for {
			if err := Seek(start, token.AccessToken); err != nil {
				fmt.Println("seek operation failed")
				return
			}

			duration := end - start
			time.Sleep(time.Duration(duration) * time.Millisecond)
		}
	}()

	return nil
}
