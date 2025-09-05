package api

import (
	"net/http"

	"github.com/dzakaeryan20/dealls-hris/internal/api/handler"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/auth"
	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authService auth.Service,
	employeeService employee.Service,
	jwtSecret string,

) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Heartbeat("/health"))

	authHandler := handler.NewAuthHandler(authService)

	r.Post("/api/v1/auth/login", authHandler.Login)

	return r
}
