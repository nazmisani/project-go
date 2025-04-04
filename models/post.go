package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model       // Menambahkan ID, CreatedAt, UpdatedAt, DeletedAt
	Title     string `gorm:"not null" json:"title"`
	Body      string `gorm:"not null" json:"body"`
	UserID    uint   `json:"user_id"`
	User      User   `json:"user" gorm:"foreignKey:UserID"`
}