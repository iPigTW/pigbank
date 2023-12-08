package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"unique" gorm:"not null"`
	Email string `gorm:"unique" gorm:"not null"`
	Password string `gorm:"not null"`
	Balance int `gorm:"default:0"`
}