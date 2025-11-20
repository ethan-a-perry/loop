package spotify

import (
	"fmt"
	"net/http"
)

func Seek(start int, accessToken string) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", start)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer " + accessToken)

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
