package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"barang-api/internal/models"
	"barang-api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BarangHandler struct {
	DB *gorm.DB
}

func NewBarangHandler(db *gorm.DB) *BarangHandler {
	return &BarangHandler{DB: db}
}

type barangRequest struct {
	SKU       string `json:"sku" binding:"required,max=50"`
	Nama      string `json:"nama" binding:"required,max=150"`
	Deskripsi string `json:"deskripsi" binding:"max=5000"`
	Harga     *int64 `json:"harga" binding:"required,gte=0"`
	Stok      *int64 `json:"stok" binding:"required,gte=0"`
}

func (h *BarangHandler) List(c *gin.Context) {
	page, err := parsePositiveInt(c.DefaultQuery("page", "1"), 1, 1_000_000)
	if err != nil {
		utils.Failure(c, http.StatusBadRequest, "Parameter page tidak valid")
		return
	}

	limit, err := parsePositiveInt(c.DefaultQuery("limit", "10"), 10, 100)
	if err != nil {
		utils.Failure(c, http.StatusBadRequest, "Parameter limit harus antara 1 dan 100")
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	if len(search) > 100 {
		utils.Failure(c, http.StatusBadRequest, "Panjang search maksimal 100 karakter")
		return
	}

	sortBy := strings.ToLower(strings.TrimSpace(c.DefaultQuery("sort_by", "created_at")))
	sortOrder := strings.ToLower(strings.TrimSpace(c.DefaultQuery("sort_order", "desc")))

	allowedSort := map[string]string{
		"id":         "id",
		"sku":        "sku",
		"nama":       "nama",
		"harga":      "harga",
		"stok":       "stok",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	column, ok := allowedSort[sortBy]
	if !ok {
		utils.Failure(c, http.StatusBadRequest, "sort_by tidak didukung")
		return
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		utils.Failure(c, http.StatusBadRequest, "sort_order harus asc atau desc")
		return
	}

	query := h.DB.Model(&models.Barang{})

	if search != "" {
		pattern := "%" + escapeLike(search) + "%"
		query = query.Where(
			"sku LIKE ? ESCAPE '!' OR nama LIKE ? ESCAPE '!'",
			pattern,
			pattern,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Failure(c, http.StatusInternalServerError, "Gagal menghitung data barang")
		return
	}

	offset := (page - 1) * limit

	var barang []models.Barang
	if err := query.
		Order(fmt.Sprintf("`%s` %s", column, strings.ToUpper(sortOrder))).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&barang).Error; err != nil {
		utils.Failure(c, http.StatusInternalServerError, "Gagal mengambil data barang")
		return
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	utils.Success(c, http.StatusOK, "Data barang berhasil diambil", gin.H{
		"items": barang,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total_items": total,
			"total_pages": totalPages,
		},
		"filters": gin.H{
			"search":     search,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

func (h *BarangHandler) Show(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var barang models.Barang
	if err := h.DB.First(&barang, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Failure(c, http.StatusNotFound, "Barang tidak ditemukan")
			return
		}

		utils.Failure(c, http.StatusInternalServerError, "Gagal mengambil data barang")
		return
	}

	utils.Success(c, http.StatusOK, "Detail barang berhasil diambil", barang)
}

func (h *BarangHandler) Create(c *gin.Context) {
	var req barangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Failure(c, http.StatusBadRequest, "Input tidak valid. Periksa SKU, nama, deskripsi, harga, dan stok")
		return
	}

	req.SKU = strings.TrimSpace(req.SKU)
	req.Nama = strings.TrimSpace(req.Nama)
	req.Deskripsi = strings.TrimSpace(req.Deskripsi)

	if req.SKU == "" || req.Nama == "" {
		utils.Failure(c, http.StatusBadRequest, "SKU dan nama tidak boleh kosong")
		return
	}

	barang := models.Barang{
		SKU:       req.SKU,
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		Harga:     *req.Harga,
		Stok:      *req.Stok,
	}

	if err := h.DB.Create(&barang).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.Failure(c, http.StatusConflict, "SKU sudah digunakan")
			return
		}

		utils.Failure(c, http.StatusInternalServerError, "Gagal menambahkan barang")
		return
	}

	utils.Success(c, http.StatusCreated, "Barang berhasil ditambahkan", barang)
}

func (h *BarangHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req barangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Failure(c, http.StatusBadRequest, "Input tidak valid. Periksa SKU, nama, deskripsi, harga, dan stok")
		return
	}

	req.SKU = strings.TrimSpace(req.SKU)
	req.Nama = strings.TrimSpace(req.Nama)
	req.Deskripsi = strings.TrimSpace(req.Deskripsi)

	if req.SKU == "" || req.Nama == "" {
		utils.Failure(c, http.StatusBadRequest, "SKU dan nama tidak boleh kosong")
		return
	}

	var barang models.Barang
	if err := h.DB.First(&barang, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Failure(c, http.StatusNotFound, "Barang tidak ditemukan")
			return
		}

		utils.Failure(c, http.StatusInternalServerError, "Gagal mengambil data barang")
		return
	}

	barang.SKU = req.SKU
	barang.Nama = req.Nama
	barang.Deskripsi = req.Deskripsi
	barang.Harga = *req.Harga
	barang.Stok = *req.Stok

	if err := h.DB.Save(&barang).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.Failure(c, http.StatusConflict, "SKU sudah digunakan oleh barang lain")
			return
		}

		utils.Failure(c, http.StatusInternalServerError, "Gagal memperbarui barang")
		return
	}

	utils.Success(c, http.StatusOK, "Barang berhasil diperbarui", barang)
}

func (h *BarangHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	result := h.DB.Delete(&models.Barang{}, id)
	if result.Error != nil {
		utils.Failure(c, http.StatusInternalServerError, "Gagal menghapus barang")
		return
	}

	if result.RowsAffected == 0 {
		utils.Failure(c, http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "Barang berhasil dihapus", gin.H{
		"id": id,
	})
}

func parsePositiveInt(value string, fallback, max int) (int, error) {
	if value == "" {
		return fallback, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil || n < 1 || n > max {
		return 0, errors.New("nilai harus berupa bilangan positif dalam batas yang diizinkan")
	}

	return n, nil
}

func parseID(c *gin.Context) (uint64, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		utils.Failure(c, http.StatusBadRequest, "ID barang tidak valid")
		return 0, false
	}

	return id, true
}

// Mencegah karakter %, _ dan ! dari input search diperlakukan sebagai wildcard SQL.
func escapeLike(value string) string {
	value = strings.ReplaceAll(value, "!", "!!")
	value = strings.ReplaceAll(value, "%", "!%")
	value = strings.ReplaceAll(value, "_", "!_")
	return value
}
