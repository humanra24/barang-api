package handlers

import (
	"net/http"
	"strings"

	"barang-api/internal/config"
	"barang-api/internal/middleware"
	"barang-api/internal/models"
	"barang-api/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type loginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=1,max=72"`
}

func NewAuthHandler(
	db *gorm.DB,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		DB:  db,
		Cfg: cfg,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Failure(
			c,
			http.StatusBadRequest,
			"Username atau password tidak valid",
		)
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	var user models.User
	err := h.DB.Where("username = ?", req.Username).
		First(&user).Error

	if err != nil {
		// Pesan yang sama untuk username tidak ditemukan
		// maupun password yang salah.
		utils.Failure(
			c,
			http.StatusUnauthorized,
			"Username atau password salah",
		)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		utils.Failure(
			c,
			http.StatusUnauthorized,
			"Username atau password salah",
		)
		return
	}

	token, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		h.Cfg.JWTSecret,
		h.Cfg.JWTExpireMinutes,
	)
	if err != nil {
		utils.Failure(
			c,
			http.StatusInternalServerError,
			"Gagal membuat token",
		)
		return
	}

	utils.Success(
		c,
		http.StatusOK,
		"Login berhasil",
		gin.H{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   h.Cfg.JWTExpireMinutes * 60,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Failure(
			c,
			http.StatusUnauthorized,
			"Autentikasi diperlukan",
		)
		return
	}

	id, ok := userID.(uint64)
	if !ok || id == 0 {
		utils.Failure(
			c,
			http.StatusUnauthorized,
			"Identitas pengguna tidak valid",
		)
		return
	}

	var user models.User

	if err := h.DB.First(&user, id).Error; err != nil {
		utils.Failure(
			c,
			http.StatusUnauthorized,
			"Akun tidak ditemukan",
		)
		return
	}

	utils.Success(
		c,
		http.StatusOK,
		"Profil pengguna",
		gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	)
}
