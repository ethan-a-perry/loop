package spotifyauth

import (
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler {
		service: service,
	}
}

func (h *handler) Connect(w http.ResponseWriter, r *http.Request) {
	authorizationUrl, err := h.service.GetAuthorizationUrl()
	if err != nil {
		http.Error(w, "Failed to generate Spotify authoirzation URL", http.StatusInternalServerError)
	}

	http.Redirect(w, r, authorizationUrl, http.StatusFound)
}

func (h *handler) Callback(w http.ResponseWriter, r *http.Request) {
	err := r.URL.Query().Get("error")
	if err != "" {
		http.Error(w, "Authorization failed during callback: " + err, http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
        return
	}

	if err := h.service.EstablishToken(code); err != nil {
		http.Error(w, "Failed to retrieve access token from Spotify", http.StatusBadRequest)
	}

	http.Redirect(w, r, "/?spotify=connected", http.StatusFound)
}
