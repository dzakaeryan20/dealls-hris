package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

// ipAddressKey adalah kunci untuk menyimpan dan mengambil IP address dari context.
const ipAddressKey contextKey = "ipAddress"

// IPTrackerMiddleware adalah middleware yang akan dijalankan untuk setiap request.
// Tugasnya adalah mendapatkan IP address klien dan menyimpannya ke dalam context.
func IPTrackerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dapatkan IP address dari request menggunakan fungsi helper.
		ip := getIPAddress(r)

		// Simpan IP address ke dalam context dari request tersebut.
		ctx := context.WithValue(r.Context(), ipAddressKey, ip)

		// Lanjutkan request ke middleware atau handler selanjutnya
		// dengan membawa context yang sudah berisi IP address.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetIPAddressFromContext adalah fungsi helper yang bisa dipanggil dari mana saja
// (misalnya dari service atau handler) untuk mengambil IP yang sudah disimpan di context.
func GetIPAddressFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(ipAddressKey).(string); ok {
		return ip
	}
	return ""
}

// getIPAddress berfungsi untuk mengambil IP address asli dari klien.
// Fungsi ini cukup pintar untuk menangani kasus di mana server berada di belakang
// sebuah reverse proxy (seperti Nginx, Caddy, atau Load Balancer).
func getIPAddress(r *http.Request) string {
	// 1. Cek header 'X-Forwarded-For', yang paling umum digunakan oleh proxy.
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Header ini bisa berisi rantai IP (client, proxy1, proxy2). IP asli klien adalah yang pertama.
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// 2. Cek header 'X-Real-IP' sebagai alternatif.
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 3. Jika tidak ada header proxy, gunakan alamat remote dari koneksi TCP sebagai fallback.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Jika ada error saat parsing (misalnya tidak ada port), kembalikan alamat aslinya.
		return r.RemoteAddr
	}
	return ip
}
