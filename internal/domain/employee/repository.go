package employee

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *Employee) error
	GetByUsername(ctx context.Context, username string) (*Employee, error)
	GetAllEmployees(ctx context.Context) ([]Employee, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, user *Employee) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*Employee, error) {
	var user Employee
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetAllEmployees(ctx context.Context) ([]Employee, error) {
	var users []Employee
	if err := r.db.WithContext(ctx).Where("role = ?", "employee").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
