package container

import (
	"lab1/models"
	"lab1/repository"
	"lab1/validation"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Container struct {
	DB               *gorm.DB
	BookRepository   repository.BookRepository
	ReaderRepository repository.ReaderRepository
	Validator        *validation.Validator
}

func NewContainer(dbPath string) (*Container, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Book{}, &models.Reader{})
	if err != nil {
		return nil, err
	}

	bookRepo := repository.NewBookRepository(db)
	readerRepo := repository.NewReaderRepository(db)

	validator := validation.NewValidator()

	return &Container{
		DB:               db,
		BookRepository:   bookRepo,
		ReaderRepository: readerRepo,
		Validator:        validator,
	}, nil
}

func (c *Container) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
