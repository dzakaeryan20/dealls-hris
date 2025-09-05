package reimbursement

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateReimbursement(ctx context.Context, reimbursement *Reimbursement) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateReimbursement(ctx context.Context, reimbursement *Reimbursement) error {
	return r.db.WithContext(ctx).Create(reimbursement).Error
}
