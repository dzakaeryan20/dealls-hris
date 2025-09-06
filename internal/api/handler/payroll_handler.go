// File: internal/api/handler/payroll_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/api/middleware"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/payroll"
	"github.com/go-chi/chi/v5"
)

// PayrollHandler menangani semua request HTTP yang berkaitan dengan penggajian.
type PayrollHandler struct {
	service payroll.Service
}

// NewPayrollHandler membuat instance baru dari PayrollHandler.
func NewPayrollHandler(s payroll.Service) *PayrollHandler {
	return &PayrollHandler{service: s}
}

// createPeriodRequest adalah struct untuk menampung data JSON saat admin membuat periode payroll.
type createPeriodRequest struct {
	StartDate string `json:"start_date" validate:"required"` // Format: "YYYY-MM-DD"
	EndDate   string `json:"end_date" validate:"required"`   // Format: "YYYY-MM-DD"
}

// CreatePayrollPeriod adalah handler untuk endpoint POST /api/v1/admin/payroll-period.
// Fungsi ini hanya bisa diakses oleh admin.
func (h *PayrollHandler) CreatePayrollPeriod(w http.ResponseWriter, r *http.Request) {
	var req createPeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Mengambil ID admin dari context untuk melacak siapa yang membuat periode ini.
	adminID := r.Context().Value(middleware.UserIDKey).(string)

	// Memanggil service untuk membuat periode payroll.
	period, err := h.service.CreatePayrollPeriod(r.Context(), startDate, endDate, adminID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(period)
}

// RunPayroll adalah handler untuk endpoint POST /api/v1/admin/payroll/{period_id}/run.
// Fungsi ini hanya bisa diakses oleh admin.
func (h *PayrollHandler) RunPayroll(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID periode dari parameter URL.
	periodID := chi.URLParam(r, "period_id")
	if periodID == "" {
		http.Error(w, "Period ID is required", http.StatusBadRequest)
		return
	}

	// Mengambil ID admin dari context untuk melacak siapa yang menjalankan payroll.
	adminID := r.Context().Value(middleware.UserIDKey).(string)

	// Memanggil service untuk menjalankan proses kalkulasi payroll.
	err := h.service.RunPayroll(r.Context(), periodID, adminID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Payroll run successfully"})
}

// GetMyPayslip adalah handler untuk endpoint GET /api/v1/payslip/{period_id}.
// Fungsi ini hanya bisa diakses oleh karyawan untuk melihat payslip mereka sendiri.
func (h *PayrollHandler) GetMyPayslip(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID pengguna (karyawan) yang sedang login dari context.
	userID := r.Context().Value(middleware.UserIDKey).(string)
	// Mengambil ID periode dari parameter URL.
	periodID := chi.URLParam(r, "period_id")

	// Memanggil service untuk mendapatkan data payslip.
	payslip, err := h.service.GetPayslip(r.Context(), userID, periodID)
	if err != nil {
		http.Error(w, "Payslip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payslip)
}

// GetPayrollSummary adalah handler untuk endpoint GET /api/v1/admin/payroll/{period_id}/summary.
// Fungsi ini hanya bisa diakses oleh admin.
func (h *PayrollHandler) GetPayrollSummary(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID periode dari parameter URL.
	periodID := chi.URLParam(r, "period_id")

	// Memanggil service untuk mendapatkan ringkasan data payroll.
	summary, err := h.service.GetPayrollSummary(r.Context(), periodID)
	if err != nil {
		http.Error(w, "Could not generate summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
