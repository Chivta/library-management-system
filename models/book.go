package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	UserID      uint   `gorm:"not null"` // Owner of the book
	User        User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
