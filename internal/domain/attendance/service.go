package attendance

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	SubmitAttendance(ctx context.Context, userID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) SubmitAttendance(ctx context.Context, userID string) error {
	today := time.Now()

	// Users cannot submit on weekends
	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		return errors.New("cannot submit attendance on a weekend")
	}

	// Submissions on the same day should count as one
	dateStr := today.Format("2006-01-02")
	hasSubmitted, err := s.repo.HasAttendanceOnDate(ctx, userID, dateStr)
	if err != nil {
		return err
	}
	if hasSubmitted {
		return errors.New("attendance already submitted for today")
	}

	attendance := &Attendance{
		UserID: userID,
		Date:   today,
	}

	return s.repo.CreateAttendance(ctx, attendance)
}
