package api

import (
	"net/http"

	"github.com/dzakaeryan20/dealls-hris/internal/api/handler"
	"github.com/dzakaeryan20/dealls-hris/internal/api/middleware"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/attendance"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/auth"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/overtime"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/payroll"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/reimbursement"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authService auth.Service,
	employeeService employee.Service,
	attendanceService attendance.Service,
	overtimeService overtime.Service,
	reimbursementService reimbursement.Service,
	payrollService payroll.Service,
	jwtSecret string,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.IPTrackerMiddleware)
	r.Use(chimiddleware.Heartbeat("/health"))

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	overtimeHandler := handler.NewOvertimeHandler(overtimeService)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementService)
	payrollHandler := handler.NewPayrollHandler(payrollService)

	// Public routes
	r.Post("/api/v1/auth/login", authHandler.Login)

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtSecret))

		// --- Employee Routes ---
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("employee"))

			// Submissions
			r.Post("/api/v1/attendance", attendanceHandler.SubmitAttendance)
			r.Post("/api/v1/overtime", overtimeHandler.SubmitOvertime)
			r.Post("/api/v1/reimbursement", reimbursementHandler.SubmitReimbursement)

			// Payslip
			r.Get("/api/v1/payslip/{period_id}", payrollHandler.GetMyPayslip)
		})

		// --- Admin Routes ---
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))

			// Payroll Management
			r.Post("/api/v1/admin/payroll-period", payrollHandler.CreatePayrollPeriod)
			r.Post("/api/v1/admin/payroll/{period_id}/run", payrollHandler.RunPayroll)
			r.Get("/api/v1/admin/payroll/{period_id}/summary", payrollHandler.GetPayrollSummary)
		})
	})

	return r
}
