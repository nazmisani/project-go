package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;index" json:"username" binding:"required"`
	Email    string `gorm:"unique;not null;index" json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=8"`
	Role     string `gorm:"default:'user'" json:"role"`
	Posts    []Post `json:"posts,omitempty" gorm:"foreignKey:UserID"`
}

