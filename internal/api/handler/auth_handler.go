package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/auth"
)

type AuthHandler struct {
	service auth.Service
}

func NewAuthHandler(s auth.Service) *AuthHandler {
	return &AuthHandler{service: s}
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Teruskan context dari request ke service
	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
