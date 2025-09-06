package reimbursement

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	SubmitReimbursement(ctx context.Context, userID string, date time.Time, description string, amount float64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) SubmitReimbursement(ctx context.Context, userID string, date time.Time, description string, amount float64) error {
	if amount <= 0 {
		return errors.New("reimbursement amount must be positive")
	}
	if description == "" {
		return errors.New("reimbursement description is required")
	}

	reimbursement := &Reimbursement{
		UserID:      userID,
		Date:        date,
		Description: description,
		Amount:      amount,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	return s.repo.CreateReimbursement(ctx, reimbursement)
}
