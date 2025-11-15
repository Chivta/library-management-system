package repository

import (
	"lab1/models"

	"gorm.io/gorm"
)

type ReaderRepository interface {
	Create(reader *models.Reader) error
	FindAll() ([]models.Reader, error)
	FindByID(id uint) (*models.Reader, error)
	Update(reader *models.Reader) error
	Delete(id uint) error
	DeleteAll() error
}

type readerRepository struct {
	db *gorm.DB
}

func NewReaderRepository(db *gorm.DB) ReaderRepository {
	return &readerRepository{db: db}
}

func (r *readerRepository) Create(reader *models.Reader) error {
	return r.db.Create(reader).Error
}

func (r *readerRepository) FindAll() ([]models.Reader, error) {
	var readers []models.Reader
	err := r.db.Find(&readers).Error
	return readers, err
}

func (r *readerRepository) FindByID(id uint) (*models.Reader, error) {
	var reader models.Reader
	err := r.db.First(&reader, id).Error
	if err != nil {
		return nil, err
	}
	return &reader, nil
}

func (r *readerRepository) Update(reader *models.Reader) error {
	return r.db.Save(reader).Error
}

func (r *readerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reader{}, id).Error
}

func (r *readerRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Reader{}).Error
}
