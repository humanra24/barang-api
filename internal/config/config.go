
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret        string
	JWTExpireMinutes int

	AdminUsername string
	AdminPassword string
}

func Load() (*Config, error) {
	// Environment variable dari sistem memiliki prioritas
	// karena godotenv.Load tidak menimpa nilai yang sudah ada.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("gagal membaca .env: %w", err)
	}

	cfg := &Config{
		Env:        getEnv("APP_ENV", "development"),
		Port:       getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		JWTSecret:        os.Getenv("JWT_SECRET"),
		AdminUsername:    os.Getenv("ADMIN_USERNAME"),
		AdminPassword:    os.Getenv("ADMIN_PASSWORD"),
		JWTExpireMinutes: 60,
	}

	appPort, err := strconv.Atoi(cfg.Port)
	if err != nil || appPort < 1 || appPort > 65535 {
		return nil, fmt.Errorf("APP_PORT tidak valid")
	}

	dbPort, err := strconv.Atoi(cfg.DBPort)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return nil, fmt.Errorf("DB_PORT tidak valid")
	}

	expireMinutes, err := strconv.Atoi(
		getEnv("JWT_EXPIRE_MINUTES", "60"),
	)
	if err != nil || expireMinutes < 1 || expireMinutes > 10080 {
		return nil, fmt.Errorf(
			"JWT_EXPIRE_MINUTES harus antara 1 dan 10080",
		)
	}
	cfg.JWTExpireMinutes = expireMinutes

	if cfg.Env != "development" && cfg.Env != "production" {
		return nil, fmt.Errorf(
			"APP_ENV harus development atau production",
		)
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, fmt.Errorf(
			"DB_HOST, DB_USER, dan DB_NAME wajib diisi",
		)
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf(
			"JWT_SECRET wajib minimal 32 karakter",
		)
	}

	if strings.TrimSpace(cfg.AdminUsername) == "" {
		return nil, fmt.Errorf("ADMIN_USERNAME wajib diisi")
	}

	if len(cfg.AdminUsername) > 50 {
		return nil, fmt.Errorf(
			"ADMIN_USERNAME maksimal 50 karakter",
		)
	}

	if len(cfg.AdminPassword) < 12 ||
		len(cfg.AdminPassword) > 72 {
		return nil, fmt.Errorf(
			"ADMIN_PASSWORD harus sepanjang 12 sampai 72 byte",
		)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
