package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"barang-api/internal/config"
	"barang-api/internal/database"
	"barang-api/internal/handlers"
	"barang-api/internal/middleware"
	"barang-api/internal/utils"

	"github.com/gin-gonic/gin"
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
	log.Println("Migrasi database berhasil")

	if err := database.EnsureAdmin(
		db,
		cfg.AdminUsername,
		cfg.AdminPassword,
	); err != nil {
		log.Fatalf("Inisialisasi admin gagal: %v", err)
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			utils.Failure(
				c,
				http.StatusServiceUnavailable,
				"Database tidak tersedia",
			)
			return
		}

		utils.Success(c, http.StatusOK, "API dan database berjalan", gin.H{
			"status": "ok",
		})
	})

	authHandler := handlers.NewAuthHandler(db, cfg)
	barangHandler := handlers.NewBarangHandler(db)

	api := router.Group("/api")

	// Endpoint login tidak membutuhkan token.
	api.POST("/auth/login", authHandler.Login)

	// Semua endpoint di dalam grup ini membutuhkan JWT yang valid.
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(cfg.JWTSecret))

	protected.GET("/auth/me", authHandler.Me)

	protected.GET("/barang", barangHandler.List)
	protected.GET("/barang/:id", barangHandler.Show)
	protected.POST("/barang", barangHandler.Create)
	protected.PUT("/barang/:id", barangHandler.Update)
	protected.DELETE("/barang/:id", barangHandler.Delete)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("Server berjalan di http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server gagal: %v", err)
	}
}
