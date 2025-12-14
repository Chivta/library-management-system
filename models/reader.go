package models

import "gorm.io/gorm"

type Reader struct {
	gorm.Model
	Name            string `gorm:"not null"`
	Surname         string `gorm:"not null"`
	CurrentlyReading []Book `gorm:"many2many:reader_books;constraint:OnDelete:CASCADE;"` // Books currently being read
}
