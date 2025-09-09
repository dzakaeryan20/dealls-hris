package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/api/middleware"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/reimbursement"
)

type ReimbursementHandler struct {
	service reimbursement.Service
}

func NewReimbursementHandler(s reimbursement.Service) *ReimbursementHandler {
	return &ReimbursementHandler{service: s}
}

type reimbursementRequest struct {
	Date        string  `json:"date"` // "YYYY-MM-DD"
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func (h *ReimbursementHandler) SubmitReimbursement(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req reimbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	// Teruskan context dari request
	if err := h.service.SubmitReimbursement(r.Context(), userID, date, req.Description, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reimbursement submitted successfully"})
}
