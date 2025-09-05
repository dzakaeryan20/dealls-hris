package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/api/middleware"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/overtime"
)

type OvertimeHandler struct {
	service overtime.Service
}

func NewOvertimeHandler(s overtime.Service) *OvertimeHandler {
	return &OvertimeHandler{service: s}
}

type overtimeRequest struct {
	Date  string `json:"date"` // "YYYY-MM-DD"
	Hours int    `json:"hours"`
}

func (h *OvertimeHandler) SubmitOvertime(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req overtimeRequest
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
	if err := h.service.SubmitOvertime(r.Context(), userID, date, req.Hours); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Overtime submitted successfully"})
}
