package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password"`
	Posts    []Post `json:"posts" gorm:"foreignKey:UserID"`
}

