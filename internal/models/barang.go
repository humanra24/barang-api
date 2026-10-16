package models

import (
	"time"

	"gorm.io/gorm"
)

type Barang struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	SKU       string `gorm:"size:50;not null;uniqueIndex" json:"sku"`
	Nama      string `gorm:"size:150;not null;index" json:"nama"`
	Deskripsi string `gorm:"type:text" json:"deskripsi"`

	// Harga dalam satuan rupiah bulat.
	// Contoh: 15000 = Rp15.000.
	Harga int64 `gorm:"type:bigint;not null;default:0" json:"harga"`

	Stok int64 `gorm:"type:bigint;not null;default:0" json:"stok"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
