package repository

import (
	"lab1/cache"
	"lab1/models"
	"log"

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
	db    *gorm.DB
	cache *cache.Cache
}

func NewReaderRepository(db *gorm.DB, cache *cache.Cache) ReaderRepository {
	return &readerRepository{db: db, cache: cache}
}

func (r *readerRepository) Create(reader *models.Reader) error {
	log.Printf("ReaderRepository.Create: creating reader with name='%s %s'", reader.Name, reader.Surname)
	err := r.db.Create(reader).Error
	if err != nil {
		log.Printf("ReaderRepository.Create: error creating reader: %v", err)
		return err
	}
	log.Printf("ReaderRepository.Create: reader created successfully with ID=%d", reader.ID)
	r.cache.Invalidate(cache.ReaderListKey())
	return nil
}

func (r *readerRepository) FindAll() ([]models.Reader, error) {
	log.Printf("ReaderRepository.FindAll: fetching all readers")
	if cached, found := r.cache.Get(cache.ReaderListKey()); found {
		log.Printf("ReaderRepository.FindAll: returning cached readers")
		return cached.([]models.Reader), nil
	}

	var readers []models.Reader
	err := r.db.Find(&readers).Error
	if err != nil {
		log.Printf("ReaderRepository.FindAll: error fetching readers: %v", err)
		return readers, err
	}

	log.Printf("ReaderRepository.FindAll: found %d readers from database", len(readers))
	r.cache.Set(cache.ReaderListKey(), readers)
	return readers, nil
}

func (r *readerRepository) FindByID(id uint) (*models.Reader, error) {
	log.Printf("ReaderRepository.FindByID: fetching reader with ID=%d", id)
	if cached, found := r.cache.Get(cache.ReaderIDKey(id)); found {
		log.Printf("ReaderRepository.FindByID: returning cached reader with ID=%d", id)
		reader := cached.(*models.Reader)
		return reader, nil
	}

	var reader models.Reader
	err := r.db.First(&reader, id).Error
	if err != nil {
		log.Printf("ReaderRepository.FindByID: error fetching reader with ID=%d: %v", id, err)
		return nil, err
	}

	log.Printf("ReaderRepository.FindByID: found reader with ID=%d from database", id)
	r.cache.Set(cache.ReaderIDKey(id), &reader)
	return &reader, nil
}

func (r *readerRepository) Update(reader *models.Reader) error {
	log.Printf("ReaderRepository.Update: updating reader with ID=%d, name='%s %s'", reader.ID, reader.Name, reader.Surname)
	err := r.db.Save(reader).Error
	if err != nil {
		log.Printf("ReaderRepository.Update: error updating reader with ID=%d: %v", reader.ID, err)
		return err
	}
	log.Printf("ReaderRepository.Update: reader with ID=%d updated successfully", reader.ID)
	r.cache.Invalidate(cache.ReaderIDKey(reader.ID))
	r.cache.Invalidate(cache.ReaderListKey())
	return nil
}

func (r *readerRepository) Delete(id uint) error {
	log.Printf("ReaderRepository.Delete: deleting reader with ID=%d", id)
	err := r.db.Delete(&models.Reader{}, id).Error
	if err != nil {
		log.Printf("ReaderRepository.Delete: error deleting reader with ID=%d: %v", id, err)
		return err
	}
	log.Printf("ReaderRepository.Delete: reader with ID=%d deleted successfully", id)
	r.cache.Invalidate(cache.ReaderIDKey(id))
	r.cache.Invalidate(cache.ReaderListKey())
	return nil
}

func (r *readerRepository) DeleteAll() error {
	log.Printf("ReaderRepository.DeleteAll: deleting all readers")
	err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Reader{}).Error
	if err != nil {
		log.Printf("ReaderRepository.DeleteAll: error deleting all readers: %v", err)
		return err
	}
	log.Printf("ReaderRepository.DeleteAll: all readers deleted successfully")
	r.cache.InvalidatePattern(cache.ReaderListKey())
	return nil
}
