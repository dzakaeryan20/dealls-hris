package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository adalah implementasi mock untuk user.Repository
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, user *employee.Employee) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByUsername(ctx context.Context, username string) (*employee.Employee, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*employee.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetAllEmployees(ctx context.Context) ([]employee.Employee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]employee.Employee), args.Error(1)
}

func TestAuthService(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	authService := NewService(mockEmployeeRepo, "test-secret")
	ctx := context.Background()

	t.Run("Login - Success", func(t *testing.T) {
		// Arrange
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &employee.Employee{
			ID:           "user-123",
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
			Role:         "employee",
		}

		mockEmployeeRepo.On("GetByUsername", ctx, "testuser").Return(mockUser, nil).Once()

		// Act
		token, err := authService.Login(ctx, "testuser", password)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockEmployeeRepo.AssertExpectations(t)
	})

	t.Run("Login - Fail User Not Found", func(t *testing.T) {
		// Arrange
		mockEmployeeRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Act
		token, err := authService.Login(ctx, "nonexistent", "password123")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		mockEmployeeRepo.AssertExpectations(t)
	})

	t.Run("Login - Fail Wrong Password", func(t *testing.T) {
		// Arrange
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &employee.Employee{
			ID:           "user-123",
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
			Role:         "employee",
		}
		mockEmployeeRepo.On("GetByUsername", ctx, "testuser").Return(mockUser, nil).Once()

		// Act
		token, err := authService.Login(ctx, "testuser", "wrongpassword")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		mockEmployeeRepo.AssertExpectations(t)
	})
}
