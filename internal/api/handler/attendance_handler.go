package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dzakaeryan20/dealls-hris/internal/api/middleware"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/attendance"
)

type AttendanceHandler struct {
	service attendance.Service
}

func NewAttendanceHandler(s attendance.Service) *AttendanceHandler {
	return &AttendanceHandler{service: s}
}

func (h *AttendanceHandler) SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Teruskan context dari request
	if err := h.service.SubmitAttendance(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Attendance submitted successfully for today"})
}
