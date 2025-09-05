package attendance

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAttendance(ctx context.Context, attendance *Attendance) error
	HasAttendanceOnDate(ctx context.Context, userID string, date string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateAttendance(ctx context.Context, attendance *Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

func (r *repository) HasAttendanceOnDate(ctx context.Context, userID string, date string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Attendance{}).Where("user_id = ? AND date = ?", userID, date).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
