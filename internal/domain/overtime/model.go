package overtime

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Overtime struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Date      time.Time `gorm:"type:date"`
	Hours     int
	CreatedAt time.Time
}

func (o *Overtime) BeforeCreate(tx *gorm.DB) error {
	o.ID = uuid.New().String()
	return nil
}
