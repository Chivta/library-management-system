package models

import "gorm.io/gorm"

type Reader struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Surname string `gorm:"not null"`
}
