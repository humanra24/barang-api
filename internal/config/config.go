package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() (*Config, error) {
	// Jika file .env tidak ditemukan, aplikasi tetap bisa
	// menggunakan environment variable dari sistem.
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("gagal membaca .env: %w", err)
		}
	}

	cfg := &Config{
		Env:        getEnv("APP_ENV", "development"),
		Port:       getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER wajib diisi")
	}

	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME wajib diisi")
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("APP_PORT tidak valid")
	}

	dbPort, err := strconv.Atoi(cfg.DBPort)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return nil, fmt.Errorf("DB_PORT tidak valid")
	}

	if cfg.Env != "development" && cfg.Env != "production" {
		return nil, fmt.Errorf(
			"APP_ENV harus development atau production",
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
