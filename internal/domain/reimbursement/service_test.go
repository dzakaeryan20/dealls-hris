package reimbursement

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubmissionRepository adalah implementasi mock untuk submission.Repository
type MockReimbursementRepository struct {
	mock.Mock
}

func (m *MockReimbursementRepository) CreateReimbursement(ctx context.Context, reimbursement *Reimbursement) error {
	args := m.Called(ctx, reimbursement)
	return args.Error(0)
}

// ... (kode yang sudah ada di service_test.go)

func TestReimbursement(t *testing.T) {
	t.Run("SubmitReimbursement - Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockReimbursementRepository)
		reimbursementService := NewService(mockRepo)
		ctx := context.Background()
		userID := "user-456"

		// Siapkan ekspektasi: Saat CreateReimbursement dipanggil dengan data apa pun, return nil (sukses).
		mockRepo.On("CreateReimbursement", ctx, mock.AnythingOfType("*submission.Reimbursement")).Return(nil).Once()

		// Act
		err := reimbursementService.SubmitReimbursement(ctx, userID, time.Now(), "Biaya Transport", 75000)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SubmitReimbursement - Fail because amount is zero or negative", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockReimbursementRepository) // mock tidak akan dipanggil
		reimbursementService := NewService(mockRepo)
		ctx := context.Background()
		userID := "user-456"

		// Act
		err := reimbursementService.SubmitReimbursement(ctx, userID, time.Now(), "Invalid Amount", -50000)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reimbursement amount must be positive")
	})

	t.Run("SubmitReimbursement - Fail because description is empty", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockReimbursementRepository) // mock tidak akan dipanggil
		reimbursementService := NewService(mockRepo)
		ctx := context.Background()
		userID := "user-456"

		// Act
		err := reimbursementService.SubmitReimbursement(ctx, userID, time.Now(), "", 75000)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reimbursement description is required")
	})
}
