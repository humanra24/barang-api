package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"barang-api/internal/config"
	"barang-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan koneksi SQL: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("database tidak dapat dihubungi: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Barang{},
	); err != nil {
		return fmt.Errorf("gagal menjalankan migrasi: %w", err)
	}

	return nil
}

func EnsureAdmin(db *gorm.DB, username, password string) error {
	var user models.User

	err := db.Where("username = ?", username).First(&user).Error
	if err == nil {
		// Admin sudah ada. Jangan menimpa password yang sudah tersimpan.
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("gagal memeriksa admin awal: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password admin: %w", err)
	}

	user = models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := db.Create(&user).Error; err != nil {
		// Menangani kemungkinan admin dibuat secara bersamaan oleh proses lain.
		var existing models.User
		if findErr := db.Where("username = ?", username).
			First(&existing).Error; findErr == nil {
			return nil
		}

		return fmt.Errorf("gagal membuat admin awal: %w", err)
	}

	log.Println("Admin awal berhasil dibuat")
	return nil
}
