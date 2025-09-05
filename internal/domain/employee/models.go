package employee

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	ID           string `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string  // 'admin' or 'employee'
	BaseSalary   float64 // Only for employees
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

func NewUser(username, password, role string, salary float64) (*Employee, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Employee{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		BaseSalary:   salary,
	}, nil
}
