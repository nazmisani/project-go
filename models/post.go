package models

type Post struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"not null"`
	Body   string `gorm:"not null"`
	UserID uint   `json:"user_id"`
}