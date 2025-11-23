package spotify

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
)

type Service struct {
	auth *spotifyauth.Service
	isLoopActive atomic.Bool
	stop chan struct{}
}

func NewService(auth *spotifyauth.Service) *Service {
	return &Service {
		auth: auth,
		stop: nil,
	}
}

func (s *Service) StartLoop(start, end int) error {
	// Wait for previous loop to stop
	if s.IsLoopActive() {
		s.StopLoop()

		timeout := time.After(5 * time.Second)
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
				case <-ticker.C:
					if !s.IsLoopActive() {
						goto done
					}
				case <-timeout:
					return fmt.Errorf("timeout waiting for the preexisting loop to stop")
			}
		}
		done:
	}

	s.setLoopActive()
	s.stop = make(chan struct{})

	go s.runLoop(start, end)

	return nil
}

func (s *Service) StopLoop() error {
	if !s.IsLoopActive() {
		return fmt.Errorf("Loop is not running")
	}

	if s.stop != nil {
		close(s.stop)
	}

	return nil
}

func (s *Service) IsLoopActive() bool {
	return s.isLoopActive.Load()
}

func (s *Service) CheckPlaybackState(accessToken, currentTrackID string) bool {
	playbackState, err := GetPlaybackState(accessToken)
	if err != nil {
		fmt.Println("failed to get playback state")
		return false
	}

	if !playbackState.Device.IsActive || !playbackState.IsPlaying {
		fmt.Println("playback not active, stopping loop")
		return false
	}

	if currentTrackID != playbackState.Item.ID {
		fmt.Println("a new track was selected, stopping loop")
		return false
	}

	return true
}

func (s *Service) runLoop(start, end int) {
	defer s.setLoopInactive()

	duration := time.Duration(end - start) * time.Millisecond
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	firstIteration := true
	currentTrackID := ""

	for {
		token, err := s.auth.GetValidToken()
		if err != nil {
			fmt.Println("failed to get valid token")
			return
		}

		if firstIteration {
			playbackState, err := GetPlaybackState(token.AccessToken)
			if err != nil {
				fmt.Println("failed to get playback state")
				return
			}
			currentTrackID = playbackState.Item.ID

			firstIteration = false
		} else if s.CheckPlaybackState(token.AccessToken, currentTrackID) == false {
			return
		}

		if err := Seek(start, token.AccessToken); err != nil {
			fmt.Println("seek operation failed")
			return
		}

		select {
			case <-ticker.C:
				continue
			case <-s.stop:
				return
		}
	}
}

func (s *Service) setLoopActive() {
	s.isLoopActive.Store(true)
}

func (s *Service) setLoopInactive() {
	s.isLoopActive.Store(false)
}
