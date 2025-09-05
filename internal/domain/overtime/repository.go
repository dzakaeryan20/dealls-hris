package overtime

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateOvertime(ctx context.Context, overtime *Overtime) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateOvertime(ctx context.Context, overtime *Overtime) error {
	return r.db.WithContext(ctx).Create(overtime).Error
}
