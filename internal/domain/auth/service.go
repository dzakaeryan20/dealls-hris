// File: internal/domain/auth/service.go
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/employee"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type service struct {
	userRepo  employee.Repository
	jwtSecret string
}

func NewService(userRepo employee.Repository, jwtSecret string) Service {
	return &service{userRepo, jwtSecret}
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	// Teruskan context ke repository
	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return generateToken(u, s.jwtSecret)
}

func generateToken(u *employee.Employee, secret string) (string, error) {
	claims := &Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
