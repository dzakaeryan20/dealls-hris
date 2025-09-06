package overtime

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	// SubmitAttendance(ctx context.Context, userID string) error
	SubmitOvertime(ctx context.Context, userID string, date time.Time, hours int) error
	// SubmitReimbursement(ctx context.Context, userID string, date time.Time, description string, amount float64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) SubmitOvertime(ctx context.Context, userID string, date time.Time, hours int) error {
	if date.Format("2006-01-02") == time.Now().Format("2006-01-02") && time.Now().Hour() < 17 {
		return errors.New("overtime can only be submitted after 5 PM")
	}

	if hours <= 0 || hours > 3 {
		return errors.New("overtime must be between 1 and 3 hours")
	}

	overtime := &Overtime{
		UserID:    userID,
		Date:      date,
		Hours:     hours,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	return s.repo.CreateOvertime(ctx, overtime)
}
