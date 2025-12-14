package repository

import (
	"lab1/cache"
	"lab1/models"
	"log"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *models.Book) error
	FindAll() ([]models.Book, error)
	FindByID(id uint) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
	DeleteAll() error
}

type bookRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewBookRepository(db *gorm.DB, cache *cache.Cache) BookRepository {
	return &bookRepository{db: db, cache: cache}
}

func (r *bookRepository) Create(book *models.Book) error {
	log.Printf("BookRepository.Create: creating book with title='%s'", book.Title)
	err := r.db.Create(book).Error
	if err != nil {
		log.Printf("BookRepository.Create: error creating book: %v", err)
		return err
	}
	log.Printf("BookRepository.Create: book created successfully with ID=%d", book.ID)
	r.cache.Invalidate(cache.BookListKey())
	return nil
}

func (r *bookRepository) FindAll() ([]models.Book, error) {
	log.Printf("BookRepository.FindAll: fetching all books")
	if cached, found := r.cache.Get(cache.BookListKey()); found {
		log.Printf("BookRepository.FindAll: returning cached books")
		return cached.([]models.Book), nil
	}

	var books []models.Book
	err := r.db.Preload("User").Find(&books).Error
	if err != nil {
		log.Printf("BookRepository.FindAll: error fetching books: %v", err)
		return books, err
	}

	log.Printf("BookRepository.FindAll: found %d books from database", len(books))
	r.cache.Set(cache.BookListKey(), books)
	return books, nil
}

func (r *bookRepository) FindByID(id uint) (*models.Book, error) {
	log.Printf("BookRepository.FindByID: fetching book with ID=%d", id)
	if cached, found := r.cache.Get(cache.BookIDKey(id)); found {
		log.Printf("BookRepository.FindByID: returning cached book with ID=%d", id)
		book := cached.(*models.Book)
		return book, nil
	}

	var book models.Book
	err := r.db.Preload("User").First(&book, id).Error
	if err != nil {
		log.Printf("BookRepository.FindByID: error fetching book with ID=%d: %v", id, err)
		return nil, err
	}

	log.Printf("BookRepository.FindByID: found book with ID=%d from database", id)
	r.cache.Set(cache.BookIDKey(id), &book)
	return &book, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	log.Printf("BookRepository.Update: updating book with ID=%d, title='%s'", book.ID, book.Title)
	err := r.db.Save(book).Error
	if err != nil {
		log.Printf("BookRepository.Update: error updating book with ID=%d: %v", book.ID, err)
		return err
	}
	log.Printf("BookRepository.Update: book with ID=%d updated successfully", book.ID)
	r.cache.Invalidate(cache.BookIDKey(book.ID))
	r.cache.Invalidate(cache.BookListKey())
	return nil
}

func (r *bookRepository) Delete(id uint) error {
	log.Printf("BookRepository.Delete: deleting book with ID=%d", id)
	err := r.db.Delete(&models.Book{}, id).Error
	if err != nil {
		log.Printf("BookRepository.Delete: error deleting book with ID=%d: %v", id, err)
		return err
	}
	log.Printf("BookRepository.Delete: book with ID=%d deleted successfully", id)
	r.cache.Invalidate(cache.BookIDKey(id))
	r.cache.Invalidate(cache.BookListKey())
	return nil
}

func (r *bookRepository) DeleteAll() error {
	log.Printf("BookRepository.DeleteAll: deleting all books")
	err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Book{}).Error
	if err != nil {
		log.Printf("BookRepository.DeleteAll: error deleting all books: %v", err)
		return err
	}
	log.Printf("BookRepository.DeleteAll: all books deleted successfully")
	r.cache.InvalidatePattern("books:") // Invalidate ALL book-related cache entries
	return nil
}
