package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"barang-api/internal/config"
	"barang-api/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Konfigurasi gagal: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database gagal: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan koneksi SQL: %v", err)
	}
	defer sqlDB.Close()

	log.Println("Koneksi MariaDB berhasil")

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migrasi database gagal: %v", err)
	}

	log.Println("Migrasi tabel berhasil")

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "Database tidak tersedia",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "API dan database berjalan",
			"data": gin.H{
				"status": "ok",
			},
		})
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Server berjalan di http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf("Server gagal: %v", err)
	}
}
