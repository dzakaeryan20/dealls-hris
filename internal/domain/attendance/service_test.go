package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubmissionRepository adalah implementasi mock untuk submission.Repository
type MockAttendanceRepository struct {
	mock.Mock
}

func (m *MockAttendanceRepository) CreateAttendance(ctx context.Context, attendance *Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepository) HasAttendanceOnDate(ctx context.Context, userID string, date string) (bool, error) {
	args := m.Called(ctx, userID, date)
	return args.Bool(0), args.Error(1)
}

func TestSubmissionService(t *testing.T) {
	t.Run("SubmitAttendance - Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockAttendanceRepository)
		submissionService := NewService(mockRepo)
		ctx := context.Background()
		userID := "user-123"

		// Asumsikan hari ini bukan weekend untuk tes ini
		todayStr := time.Now().Format("2006-01-02")

		mockRepo.On("HasAttendanceOnDate", ctx, userID, todayStr).Return(false, nil).Once()
		mockRepo.On("CreateAttendance", ctx, mock.AnythingOfType("*submission.Attendance")).Return(nil).Once()

		// Act
		err := submissionService.SubmitAttendance(ctx, userID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SubmitAttendance - Fail because already submitted", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockAttendanceRepository)
		submissionService := NewService(mockRepo)
		ctx := context.Background()
		userID := "user-123"
		todayStr := time.Now().Format("2006-01-02")

		mockRepo.On("HasAttendanceOnDate", ctx, userID, todayStr).Return(true, nil).Once()

		// Act
		err := submissionService.SubmitAttendance(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attendance already submitted")
		mockRepo.AssertExpectations(t)
	})

}
