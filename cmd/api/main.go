package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dzakaeryan20/dealls-hris/internal/api"
	"github.com/dzakaeryan20/dealls-hris/internal/config"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/attendance"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/auth"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/overtime"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/payroll"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/reimbursement"
	"github.com/dzakaeryan20/dealls-hris/internal/platform/database"
	"github.com/dzakaeryan20/dealls-hris/internal/platform/seeder"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 2. Initialize Database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&employee.Employee{},
		&payroll.PayrollPeriod{},
		&attendance.Attendance{},
		&overtime.Overtime{},
		&reimbursement.Reimbursement{},
		&payroll.Payslip{},
	)
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	// 3. Run Seeder (optional)
	if cfg.RunSeeder {
		log.Println("Running database seeder...")
		if err := seeder.Run(db); err != nil {
			log.Fatalf("seeder failed: %v", err)
		}
		log.Println("Seeder finished successfully.")
	}

	// 4. Initialize Repositories
	userRepo := employee.NewRepository(db)
	attendanceRepo := attendance.NewRepository(db)
	overtimeRepo := overtime.NewRepository(db)
	reimbursementRepo := reimbursement.NewRepository(db)
	payrollRepo := payroll.NewRepository(db)

	// 5. Initialize Services
	authService := auth.NewService(userRepo, cfg.JWTSecret)
	employeeService := employee.NewService(userRepo)
	attendanceService := attendance.NewService(attendanceRepo)
	overtimeService := overtime.NewService(overtimeRepo)
	reimbursementService := reimbursement.NewService(reimbursementRepo)
	payrollService := payroll.NewService(payrollRepo, userRepo)

	// 6. Initialize Router
	router := api.NewRouter(authService,
		employeeService,
		attendanceService,
		overtimeService,
		reimbursementService,
		payrollService, cfg.JWTSecret)

	// 7. Start Server
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
