package reimbursement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reimbursement struct {
	ID          string    `gorm:"primaryKey"`
	UserID      string    `gorm:"index"`
	Date        time.Time `gorm:"type:date"`
	Description string
	Amount      float64
	CreatedAt   time.Time
}

func (r *Reimbursement) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New().String()
	return nil
}
