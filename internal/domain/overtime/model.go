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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `gorm:"size:36" json:"created_by"`
	UpdatedBy string    `gorm:"size:36" json:"updated_by"`
}

func (o *Overtime) BeforeCreate(tx *gorm.DB) error {
	o.ID = uuid.New().String()
	return nil
}
