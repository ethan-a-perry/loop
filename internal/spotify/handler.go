package spotify

import (
	"encoding/json"
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

func (h *handler) Loop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req := struct {
		Start int `json:"start"`
		End int `json:"end"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON (request body)", http.StatusBadRequest)
		return
	}

	if err := h.service.Loop(req.Start, req.End); err != nil {
		http.Error(w, "Authorization required. Please authenticate with Spotify.", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *handler) StopLoop(w http.ResponseWriter, r *http.Request) {

}
