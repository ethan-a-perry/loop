package spotifyauth

import (
	"fmt"
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

func (h *handler) Callback(w http.ResponseWriter, r *http.Request) {
	err := r.URL.Query().Get("error")

	if err != "" {
		http.Error(w, "Authorization failed during callback" + err, http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
        return
	}

	h.service.EstablishToken(code)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *handler) RequestSpotify(w http.ResponseWriter, r *http.Request) {
	authorizationUrl, err := h.service.GetAuthorizationUrl()

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, authorizationUrl, http.StatusFound)
}
