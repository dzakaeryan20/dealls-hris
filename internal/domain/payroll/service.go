// File: internal/domain/payroll/service.go
package payroll

import (
	"context"
	"errors"
	"log"

	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
)

// Service mendefinisikan kontrak untuk logika bisnis payroll.
type Service interface {
	CreatePayrollPeriod(ctx context.Context, startDate, endDate time.Time, adminID string) (*PayrollPeriod, error)
	RunPayroll(ctx context.Context, periodID string, adminID string) error
	GetPayslip(ctx context.Context, userID, periodID string) (*Payslip, error)
	GetPayrollSummary(ctx context.Context, periodID string) (*Summary, error)
}

// service adalah implementasi dari Service interface.
// Ia bergantung pada repository untuk akses data.
type service struct {
	repo         Repository
	employeeRepo employee.Repository
}

// NewService membuat instance baru dari service payroll.
func NewService(repo Repository, employee employee.Repository) Service {
	return &service{repo, employee}
}

func (s *service) CreatePayrollPeriod(ctx context.Context, startDate, endDate time.Time, adminID string) (*PayrollPeriod, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	period := NewPayrollPeriod(startDate, endDate)
	period.CreatedBy = adminID
	period.UpdatedBy = adminID

	if err := s.repo.CreatePayrollPeriod(ctx, period); err != nil {
		return nil, err
	}
	return period, nil
}

func (s *service) GetPayslip(ctx context.Context, userID, periodID string) (*Payslip, error) {
	return s.repo.GetPayslip(ctx, userID, periodID)
}

func (s *service) GetPayrollSummary(ctx context.Context, periodID string) (*Summary, error) {
	payslips, err := s.repo.GetPayslipsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}

	summary := Summary{
		PayrollPeriodID: periodID,
		EmployeePays:    []EmployeePay{},
		TotalPayout:     0.0,
	}

	for _, p := range payslips {
		summary.EmployeePays = append(summary.EmployeePays, EmployeePay{
			UserID:      p.UserID,
			TakeHomePay: p.TotalPay,
		})
		summary.TotalPayout += p.TotalPay
	}

	return &summary, nil
}

// RunPayroll adalah fungsi inti yang mengorkestrasi seluruh proses kalkulasi gaji.
func (s *service) RunPayroll(ctx context.Context, periodID string, adminID string) error {
	// 1. Ambil data periode & validasi statusnya
	period, err := s.repo.GetPayrollPeriod(ctx, periodID)
	if err != nil {
		return err
	}
	if period.Status == "completed" {
		return errors.New("payroll for this period has already been run")
	}

	// 2. Tandai periode sebagai "processing" untuk mencegah eksekusi ganda
	if err := s.repo.UpdatePayrollPeriodStatus(ctx, period.ID, "processing", adminID); err != nil {
		return err
	}

	// 3. Ambil semua data karyawan
	employees, err := s.employeeRepo.GetAllEmployees(ctx)
	if err != nil {
		s.repo.UpdatePayrollPeriodStatus(context.Background(), period.ID, "pending", adminID) // Rollback
		return err
	}

	workingDays := calculateWorkingDays(period.StartDate, period.EndDate)
	if workingDays == 0 {
		s.repo.UpdatePayrollPeriodStatus(ctx, period.ID, "completed", adminID)
		log.Println("No working days in the period. Payroll marked as completed.")
		return nil
	}

	// 4. Lakukan iterasi untuk setiap karyawan untuk menghitung gaji
	for _, emp := range employees {
		select {
		case <-ctx.Done(): // Cek apakah request dibatalkan oleh klien
			log.Println("Payroll run cancelled.")
			s.repo.UpdatePayrollPeriodStatus(context.Background(), period.ID, "pending", adminID) // Rollback
			return ctx.Err()
		default:
		}

		// Kalkulasi rate harian dan per jam
		dailyRate := emp.BaseSalary / float64(workingDays)
		hourlyRate := dailyRate / 8.0

		// Ambil semua data relevan dari repository
		attendances, _ := s.repo.GetAttendances(ctx, emp.ID, period.StartDate, period.EndDate)
		overtimes, _ := s.repo.GetOvertimes(ctx, emp.ID, period.StartDate, period.EndDate)
		reimbursements, _ := s.repo.GetReimbursements(ctx, emp.ID, period.StartDate, period.EndDate)

		// Lakukan kalkulasi
		proratedSalary := dailyRate * float64(len(attendances))
		var overtimePay float64
		for _, ot := range overtimes {
			overtimePay += float64(ot.Hours) * hourlyRate * 2
		}
		var reimbursementTotal float64
		for _, r := range reimbursements {
			reimbursementTotal += r.Amount
		}
		totalPay := proratedSalary + overtimePay + reimbursementTotal

		// Buat record payslip
		payslip := &Payslip{
			UserID:             emp.ID,
			PayrollPeriodID:    period.ID,
			BaseSalary:         emp.BaseSalary,
			ProratedSalary:     proratedSalary,
			OvertimePay:        overtimePay,
			ReimbursementTotal: reimbursementTotal,
			TotalPay:           totalPay,
			CreatedBy:          adminID,
			UpdatedBy:          adminID,
		}

		if err := s.repo.CreatePayslip(ctx, payslip); err != nil {
			log.Printf("Failed to create payslip for user %s: %v", emp.ID, err)
			continue
		}
	}

	// 5. Tandai periode sebagai "completed"
	return s.repo.UpdatePayrollPeriodStatus(ctx, period.ID, "completed", adminID)
}

func calculateWorkingDays(start, end time.Time) int {
	days := 0
	current := start
	for !current.After(end) {
		wd := current.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}
	return days
}
