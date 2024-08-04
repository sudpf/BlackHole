package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:32;unique;not null"`
	Password string `json:"password" gorm:"size:255;not null"`
	Email    string `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Phone    string `json:"phone" gorm:"size:20"`
}
