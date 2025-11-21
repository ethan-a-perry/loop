package spotify

import (
	"fmt"
	"time"

	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
)

type Service struct {
	auth *spotifyauth.Service
	stop chan struct{}
	loopActive bool
}

func NewService(auth *spotifyauth.Service) *Service {
	return &Service {
		auth: auth,
		stop: nil,
	}
}

func (s *Service) Loop(start, end int) error {
	if s.loopActive {
		close(s.stop)
	}

	token, err := s.auth.GetValidToken()
	if err != nil || token == nil {
		return fmt.Errorf("failed to get valid token")
	}

	s.stop = make(chan struct{})
	s.loopActive = true

	state, err := GetPlaybackState(token.AccessToken)
	if err != nil {
	    return fmt.Errorf("failed to get playback state: %w", err)
	}

	currentTrackID := state.Item.ID

	go func() {
		defer func() {
			s.loopActive = false
		}()

		for {
			if currentTrackID != "" {
				playbackState, err := GetPlaybackState(token.AccessToken)
				if err != nil {
					fmt.Println("failed to get playback state")
					return
				}

				if !playbackState.Device.IsActive || !playbackState.IsPlaying {
					fmt.Println("playback not active, stopping loop")
					return
				}

				if currentTrackID != playbackState.Item.ID {
					fmt.Println("a new track was selected, stopping loop")
					return
				}
			}

			if err := Seek(start, token.AccessToken); err != nil {
				fmt.Println("seek operation failed")
				return
			}

			duration := time.Duration(end - start) * time.Millisecond
			endTime := time.Now().Add(duration)

			ticker := time.NewTicker(500 * time.Millisecond)

			for time.Now().Before(endTime) {
				select {
					case <-ticker.C:
					case <-s.stop:
						ticker.Stop()
						fmt.Println("loop stopped")
						s.stop = nil
						return
				}
			}

			ticker.Stop()
		}
	}()

	return nil
}
