package attendance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"index;uniqueIndex:idx_user_date"`
	Date      time.Time `gorm:"type:date;uniqueIndex:idx_user_date"`
	CreatedAt time.Time
}

func (a *Attendance) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New().String()
	return nil
}
