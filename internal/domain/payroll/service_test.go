package payroll

import (
	"context"
	"testing"
	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/attendance"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/auth"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/overtime"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/reimbursement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPayrollRepository adalah implementasi mock untuk payroll.Repository
type MockPayrollRepository struct {
	mock.Mock
}

func (m *MockPayrollRepository) CreatePayrollPeriod(ctx context.Context, period *PayrollPeriod) error {
	args := m.Called(ctx, period)
	return args.Error(0)
}
func (m *MockPayrollRepository) GetPayrollPeriod(ctx context.Context, id string) (*PayrollPeriod, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*PayrollPeriod), args.Error(1)
}
func (m *MockPayrollRepository) UpdatePayrollPeriodStatus(ctx context.Context, id, status string, updatedByID string) error {
	args := m.Called(ctx, id, status, updatedByID)
	return args.Error(0)
}
func (m *MockPayrollRepository) GetAttendances(ctx context.Context, userID string, start, end time.Time) ([]attendance.Attendance, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]attendance.Attendance), args.Error(1)
}
func (m *MockPayrollRepository) GetOvertimes(ctx context.Context, userID string, start, end time.Time) ([]overtime.Overtime, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]overtime.Overtime), args.Error(1)
}
func (m *MockPayrollRepository) GetReimbursements(ctx context.Context, userID string, start, end time.Time) ([]reimbursement.Reimbursement, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]reimbursement.Reimbursement), args.Error(1)
}
func (m *MockPayrollRepository) CreatePayslip(ctx context.Context, payslip *Payslip) error {
	args := m.Called(ctx, payslip)
	return args.Error(0)
}
func (m *MockPayrollRepository) GetPayslip(ctx context.Context, userID, periodID string) (*Payslip, error) {
	args := m.Called(ctx, userID, periodID)
	return args.Get(0).(*Payslip), args.Error(1)
}
func (m *MockPayrollRepository) GetPayslipsByPeriod(ctx context.Context, periodID string) ([]Payslip, error) {
	args := m.Called(ctx, periodID)
	return args.Get(0).([]Payslip), args.Error(1)
}

func TestPayrollService(t *testing.T) {
	t.Run("RunPayroll - Success", func(t *testing.T) {
		// Arrange
		mockPayrollRepo := new(MockPayrollRepository)
		mockEmployeeRepo := new(auth.MockEmployeeRepository) // Menggunakan mock dari auth test
		payrollService := NewService(mockPayrollRepo, mockEmployeeRepo)

		ctx := context.Background()
		periodID := "period-001"
		adminID := "admin-001"
		startDate, _ := time.Parse("2006-01-02", "2025-09-01")
		endDate, _ := time.Parse("2006-01-02", "2025-09-05") // 5 hari, 5 hari kerja

		mockPeriod := &PayrollPeriod{
			ID:        periodID,
			StartDate: startDate,
			EndDate:   endDate,
			Status:    "pending",
		}
		mockEmployees := []employee.Employee{
			{ID: "user-001", BaseSalary: 5000000}, // Gaji 5jt, per hari 1jt
		}

		mockAttendances := []attendance.Attendance{{}, {}, {}, {}}           // 4 hari hadir
		mockOvertimes := []overtime.Overtime{{Hours: 2}}                     // 2 jam lembur
		mockReimbursements := []reimbursement.Reimbursement{{Amount: 50000}} // reimburse 50rb

		// Menyiapkan ekspektasi panggilan mock
		mockPayrollRepo.On("GetPayrollPeriod", ctx, periodID).Return(mockPeriod, nil).Once()
		mockPayrollRepo.On("UpdatePayrollPeriodStatus", ctx, periodID, "processing", adminID).Return(nil).Once()
		mockEmployeeRepo.On("GetAllEmployees", ctx).Return(mockEmployees, nil).Once()
		mockPayrollRepo.On("GetAttendances", ctx, "user-001", startDate, endDate).Return(mockAttendances, nil).Once()
		mockPayrollRepo.On("GetOvertimes", ctx, "user-001", startDate, endDate).Return(mockOvertimes, nil).Once()
		mockPayrollRepo.On("GetReimbursements", ctx, "user-001", startDate, endDate).Return(mockReimbursements, nil).Once()

		// Ekspektasi kalkulasi
		// Prorated: 1jt/hari * 4 hari = 4jt
		// Overtime: (1jt/8jam) * 2jam * 2kali = 500rb
		// Reimburse: 50rb
		// Total: 4.550.000
		mockPayrollRepo.On("CreatePayslip", ctx, mock.MatchedBy(func(p *Payslip) bool {
			return p.UserID == "user-001" && p.TotalPay == 4550000
		})).Return(nil).Once()

		mockPayrollRepo.On("UpdatePayrollPeriodStatus", ctx, periodID, "completed", adminID).Return(nil).Once()

		// Act
		err := payrollService.RunPayroll(ctx, periodID, adminID)

		// Assert
		assert.NoError(t, err)
		mockPayrollRepo.AssertExpectations(t)
		mockEmployeeRepo.AssertExpectations(t)
	})
}
