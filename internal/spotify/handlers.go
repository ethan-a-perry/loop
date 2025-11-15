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
	// Allow requests from your frontend
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if r.Method == http.MethodOptions {
		// Preflight request, respond 200 OK
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := struct {
		Start int `json:"start"`
		End int `json:"end"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON (request body)", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})

	h.service.Loop(req.Start, req.End)
}
