package payroll

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollPeriod struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status" gorm:"default:'pending'"` // pending, processing, completed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `gorm:"size:36" json:"created_by"`
	UpdatedBy string    `gorm:"size:36" json:"updated_by"`
}

func (p *PayrollPeriod) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()
	return nil
}

func NewPayrollPeriod(start, end time.Time) *PayrollPeriod {
	return &PayrollPeriod{
		StartDate: start,
		EndDate:   end,
		Status:    "pending",
	}
}

type Payslip struct {
	ID                 string    `json:"id" gorm:"primaryKey"`
	UserID             string    `json:"user_id" gorm:"index"`
	PayrollPeriodID    string    `json:"payroll_period_id" gorm:"index"`
	BaseSalary         float64   `json:"base_salary"`
	ProratedSalary     float64   `json:"prorated_salary"`
	OvertimePay        float64   `json:"overtime_pay"`
	ReimbursementTotal float64   `json:"reimbursement_total"`
	TotalPay           float64   `json:"total_pay"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedBy          string    `gorm:"size:36" json:"created_by"`
	UpdatedBy          string    `gorm:"size:36" json:"updated_by"`
}

func (p *Payslip) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()
	return nil
}

// Summary models for API responses
type Summary struct {
	PayrollPeriodID string        `json:"payroll_period_id"`
	EmployeePays    []EmployeePay `json:"employee_pays"`
	TotalPayout     float64       `json:"total_payout"`
}

type EmployeePay struct {
	UserID      string  `json:"user_id"`
	TakeHomePay float64 `json:"take_home_pay"`
}
