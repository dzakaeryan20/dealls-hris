// File: internal/api/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dzakaeryan20/dealls-hris/internal/domain/auth"
)

// contextKey adalah tipe kustom untuk kunci konteks untuk menghindari tabrakan
// dengan package lain.
type contextKey string

// UserIDKey dan UserRoleKey adalah kunci yang digunakan untuk menyimpan dan mengambil
// data pengguna dari request context.
const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

// AuthMiddleware berfungsi untuk memvalidasi JWT (JSON Web Token) dari header Authorization.
// Jika token valid, informasi pengguna (ID dan Role) akan dimasukkan ke dalam
// context dari request tersebut.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ambil header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// 2. Pastikan formatnya adalah "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Could not find bearer token in Authorization header", http.StatusUnauthorized)
				return
			}

			// 3. Validasi token
			claims, err := auth.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 4. Jika valid, masukkan UserID dan Role ke dalam context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			// 5. Lanjutkan request ke handler selanjutnya dengan context yang sudah diperbarui
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware adalah lapisan keamanan kedua setelah AuthMiddleware.
// Middleware ini memeriksa apakah role pengguna yang ada di dalam context
// cocok dengan role yang dibutuhkan untuk mengakses endpoint tertentu.
func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ambil role dari context (yang sudah dimasukkan oleh AuthMiddleware)
			role, ok := r.Context().Value(UserRoleKey).(string)

			// 2. Jika role tidak ada atau tidak cocok, tolak akses
			if !ok || role != requiredRole {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			// 3. Jika cocok, lanjutkan request ke handler
			next.ServeHTTP(w, r)
		})
	}
}
