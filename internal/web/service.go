package web

import (
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/spotify"
	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
)

type Service struct {
	authService *spotifyauth.Service
	spotifyService *spotify.Service
}

type status string

const (
	StatusUnauthorized status = "unauthorized"
	StatusIdle status = "idle"
	StatusPlaying status = "playing"
	StatusLooping status = "looping"
)

type PageState struct {
	Status status
	StatusDescription string
	StatusMessage string
	Playback *spotify.PlaybackState
	Loop struct {
		Start int
		End int
	}
}

func NewService(authService *spotifyauth.Service, spotifyService *spotify.Service) *Service {
	return &Service{
		authService: authService,
		spotifyService: spotifyService,
	}
}

func (s *Service) GetState(r *http.Request) PageState {
	if r.URL.Query().Get("spotify") == "failed" {
		return PageState{
			Status: StatusUnauthorized,
			StatusDescription: "Not connected",
			StatusMessage: "Spotify connection failed. Please try again.",
		}
	}

	token, err := s.authService.GetValidToken()
	if err != nil {
		return PageState{
			Status: StatusUnauthorized,
			StatusDescription: "Not connected",
			StatusMessage: "Please connect to Spotify to continue.",
		}
	}

	playback, err := spotify.GetPlaybackState(token.AccessToken)
	if err != nil || !playback.IsPlaying {
		return PageState{
			Status: StatusIdle,
			StatusDescription: "No active playback",
		}
	}

	if s.spotifyService.IsLoopActive() {
		start, end := s.spotifyService.GetLoopRange()

		return PageState{
			Status: StatusLooping,
			StatusDescription: "Loop active",
			Playback: playback,
			Loop: struct{
				Start int;
				End int
			}{
				Start: start,
				End: end,
			},
		}
	}

	return PageState{
		Status: StatusPlaying,
		StatusDescription: "Playing",
		Playback: playback,
	}
}
