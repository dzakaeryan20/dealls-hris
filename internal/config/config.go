package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
	RunSeeder bool
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Log that .env is not found, but don't fail, as env vars might be set in the environment
	}

	runSeeder, _ := strconv.ParseBool(os.Getenv("RUN_SEEDER"))

	return &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "admin"),
		DBPass:    getEnv("DB_PASSWORD", "secret"),
		DBName:    getEnv("DB_NAME", "payroll_db"),
		JWTSecret: getEnv("JWT_SECRET", "default_secret"),
		RunSeeder: runSeeder,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
