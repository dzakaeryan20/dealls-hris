package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	// Check if admin user already exists
	var count int64
	db.Model(&employee.Employee{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		log.Println("Seeder has already been run. Skipping.")
		return nil
	}

	ctx := context.Background()
	employeeRepo := employee.NewRepository(db)

	// Create Admin
	admin, err := employee.NewUser("admin", "password123", "admin", 0)
	if err != nil {
		return fmt.Errorf("failed to create admin user model: %w", err)
	}
	if err := employeeRepo.Create(ctx, admin); err != nil {
		return fmt.Errorf("failed to save admin user: %w", err)
	}
	log.Println("Admin user created.")

	// Create Employees
	for i := 1; i <= 100; i++ {
		username := fmt.Sprintf("employee%d", i)
		salary := gofakeit.Float64Range(4000000, 10000000) // Gaji bulanan dalam Rupiah

		employee, err := employee.NewUser(username, "password123", "employee", salary)
		if err != nil {
			log.Printf("failed to create model for %s: %v", username, err)
			continue
		}
		if err := employeeRepo.Create(ctx, employee); err != nil {
			log.Printf("failed to save employee %s: %v", username, err)
			continue
		}
	}
	log.Println("100 employee users created.")

	return nil
}
