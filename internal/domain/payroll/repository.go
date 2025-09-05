package payroll

import (
	"context"
	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/attendance"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/overtime"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/reimbursement"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePayrollPeriod(ctx context.Context, period *PayrollPeriod) error
	GetPayrollPeriod(ctx context.Context, id string) (*PayrollPeriod, error)
	UpdatePayrollPeriodStatus(ctx context.Context, id, status string) error

	GetAttendances(ctx context.Context, userID string, start, end time.Time) ([]attendance.Attendance, error)
	GetOvertimes(ctx context.Context, userID string, start, end time.Time) ([]overtime.Overtime, error)
	GetReimbursements(ctx context.Context, userID string, start, end time.Time) ([]reimbursement.Reimbursement, error)

	CreatePayslip(ctx context.Context, payslip *Payslip) error
	GetPayslip(ctx context.Context, userID, periodID string) (*Payslip, error)
	GetPayslipsByPeriod(ctx context.Context, periodID string) ([]Payslip, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreatePayrollPeriod(ctx context.Context, period *PayrollPeriod) error {
	return r.db.WithContext(ctx).Create(period).Error
}

func (r *repository) GetPayrollPeriod(ctx context.Context, id string) (*PayrollPeriod, error) {
	var period PayrollPeriod
	if err := r.db.WithContext(ctx).First(&period, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *repository) UpdatePayrollPeriodStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&PayrollPeriod{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) GetAttendances(ctx context.Context, userID string, start, end time.Time) ([]attendance.Attendance, error) {
	var attendances []attendance.Attendance
	err := r.db.WithContext(ctx).Where("user_id = ? AND date >= ? AND date <= ?", userID, start, end).Find(&attendances).Error
	return attendances, err
}

func (r *repository) GetOvertimes(ctx context.Context, userID string, start, end time.Time) ([]overtime.Overtime, error) {
	var overtimes []overtime.Overtime
	err := r.db.WithContext(ctx).Where("user_id = ? AND date >= ? AND date <= ?", userID, start, end).Find(&overtimes).Error
	return overtimes, err
}

func (r *repository) GetReimbursements(ctx context.Context, userID string, start, end time.Time) ([]reimbursement.Reimbursement, error) {
	var reimbursements []reimbursement.Reimbursement
	err := r.db.WithContext(ctx).Where("user_id = ? AND date >= ? AND date <= ?", userID, start, end).Find(&reimbursements).Error
	return reimbursements, err
}

func (r *repository) CreatePayslip(ctx context.Context, payslip *Payslip) error {
	return r.db.WithContext(ctx).Create(payslip).Error
}

func (r *repository) GetPayslip(ctx context.Context, userID, periodID string) (*Payslip, error) {
	var payslip Payslip
	err := r.db.WithContext(ctx).Where("user_id = ? AND payroll_period_id = ?", userID, periodID).First(&payslip).Error
	return &payslip, err
}

func (r *repository) GetPayslipsByPeriod(ctx context.Context, periodID string) ([]Payslip, error) {
	var payslips []Payslip
	err := r.db.WithContext(ctx).Where("payroll_period_id = ?", periodID).Find(&payslips).Error
	return payslips, err
}
