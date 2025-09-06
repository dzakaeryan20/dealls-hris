package overtime

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubmissionRepository adalah implementasi mock untuk submission.Repository
type MockOvertimeRepository struct {
	mock.Mock
}

func (m *MockOvertimeRepository) CreateOvertime(ctx context.Context, overtime *Overtime) error {
	args := m.Called(ctx, overtime)
	return args.Error(0)
}

func TestSubmissionService(t *testing.T) {
	t.Run("SubmitOvertime - Fail because hours are more than 3", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockOvertimeRepository)
		submissionService := NewService(mockRepo)
		ctx := context.Background()
		userID := "user-123"

		// Act
		err := submissionService.SubmitOvertime(ctx, userID, time.Now(), 4)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 1 and 3 hours")
	})

}
