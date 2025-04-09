package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model       // Menambahkan ID, CreatedAt, UpdatedAt, DeletedAt
	Title     string `gorm:"not null" json:"title" binding:"required"`
	Body      string `gorm:"not null" json:"body" binding:"required"`
	UserID    uint   `json:"user_id"`
	User      User   `json:"user,omitempty" gorm:"foreignKey:UserID"` // tambahkan omitempty agar tidak divalidasi
}