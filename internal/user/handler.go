package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	GetProfileHandler(w http.ResponseWriter, r *http.Request)
	UpdateProfileHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	s *service
}

func (h *handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// verify HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.s.Register(context.Background(), req)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrEmailExists {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
	}

	response := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
