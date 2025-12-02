package handlers

import (
	"errors"
	"lab1/config"
	"lab1/dto"
	"lab1/models"
	"lab1/repository"
	"lab1/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BooksHandler struct {
	repo      repository.BookRepository
	validator *validation.Validator
	config    *config.Config
}

func NewBooksHandler(repo repository.BookRepository, validator *validation.Validator, config *config.Config) *BooksHandler {
	return &BooksHandler{repo: repo, validator: validator, config: config}
}

// @Summary Get all books
// @Tags books
// @Produce json
// @Success 200 {array} dto.BookResponseDTO
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/ [get]
func (h *BooksHandler) GetAll(c *gin.Context) {
	if !h.config.EnableGetBooks {
		c.JSON(http.StatusForbidden, gin.H{"error": "GET /books endpoint is disabled"})
		return
	}

	books, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}

	response := make([]dto.BookResponseDTO, len(books))
	for i, book := range books {
		response[i] = dto.BookResponseDTO{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body dto.BookCreateDTO true "Book to create"
// @Success 201 {object} dto.BookResponseDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/ [post]
func (h *BooksHandler) Create(c *gin.Context) {
	if !h.config.EnablePostBooks {
		c.JSON(http.StatusForbidden, gin.H{"error": "POST /books endpoint is disabled"})
		return
	}

	var bookDTO dto.BookCreateDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if err := h.validator.ValidateStruct(bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidationErrors(err))
		return
	}

	book := models.Book{
		Title:       bookDTO.Title,
		Description: bookDTO.Description,
	}

	if err := h.repo.Create(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	response := dto.BookResponseDTO{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary Delete all books
// @Tags books
// @Success 204
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/ [delete]
func (h *BooksHandler) DeleteAll(c *gin.Context) {
	if !h.config.EnableDeleteBooks {
		c.JSON(http.StatusForbidden, gin.H{"error": "DELETE /books endpoint is disabled"})
		return
	}

	if err := h.repo.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete books"})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get book by ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} dto.BookResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BooksHandler) GetByID(c *gin.Context) {
	if !h.config.EnableGetBooks {
		c.JSON(http.StatusForbidden, gin.H{"error": "GET /books endpoint is disabled"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	response := dto.BookResponseDTO{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Update book by ID
// @Tags books
// @Accept json
// @Param id path int true "Book ID"
// @Param book body dto.BookUpdateDTO true "Updated book data"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [put]
func (h *BooksHandler) Update(c *gin.Context) {
	if !h.config.EnablePutBooks {
		c.JSON(http.StatusForbidden, gin.H{"error": "PUT /books endpoint is disabled"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	var bookDTO dto.BookUpdateDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if err := h.validator.ValidateStruct(bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidationErrors(err))
		return
	}

	book.Title = bookDTO.Title
	book.Description = bookDTO.Description

	if err := h.repo.Update(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Delete book by ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BooksHandler) Delete(c *gin.Context) {
	if !h.config.EnableDeleteBooks {
		c.JSON(http.StatusForbidden, gin.H{"error": "DELETE /books endpoint is disabled"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	_, err = h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.Status(http.StatusNoContent)
}
