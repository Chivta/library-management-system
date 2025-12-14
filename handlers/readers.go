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

type ReadersHandler struct {
	repo      repository.ReaderRepository
	validator *validation.Validator
	config    *config.Config
}

func NewReadersHandler(repo repository.ReaderRepository, validator *validation.Validator, config *config.Config) *ReadersHandler {
	return &ReadersHandler{repo: repo, validator: validator, config: config}
}

// @Summary Get all readers
// @Tags readers
// @Produce json
// @Success 200 {array} dto.ReaderResponseDTO
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readers/ [get]
func (h *ReadersHandler) GetAll(c *gin.Context) {
	if !h.config.EnableGetReaders {
		c.JSON(http.StatusForbidden, gin.H{"error": "GET /readers endpoint is disabled"})
		return
	}

	readers, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve readers"})
		return
	}

	response := make([]dto.ReaderResponseDTO, len(readers))
	for i, reader := range readers {
		books := make([]dto.BookResponseDTO, len(reader.CurrentlyReading))
		for j, book := range reader.CurrentlyReading {
			books[j] = dto.BookResponseDTO{
				ID:          book.ID,
				Title:       book.Title,
				Description: book.Description,
				UserID:      book.UserID,
				Username:    book.User.Username,
			}
		}
		response[i] = dto.ReaderResponseDTO{
			ID:               reader.ID,
			Name:             reader.Name,
			Surname:          reader.Surname,
			CurrentlyReading: books,
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Create a new reader
// @Tags readers
// @Accept json
// @Produce json
// @Param reader body dto.ReaderCreateDTO true "Reader to create"
// @Success 201 {object} dto.ReaderResponseDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readers/ [post]
func (h *ReadersHandler) Create(c *gin.Context) {
	if !h.config.EnablePostReaders {
		c.JSON(http.StatusForbidden, gin.H{"error": "POST /readers endpoint is disabled"})
		return
	}

	var readerDTO dto.ReaderCreateDTO
	if err := c.ShouldBindJSON(&readerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if err := h.validator.ValidateStruct(readerDTO); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidationErrors(err))
		return
	}

	reader := models.Reader{
		Name:    readerDTO.Name,
		Surname: readerDTO.Surname,
	}

	if err := h.repo.Create(&reader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reader"})
		return
	}

	response := dto.ReaderResponseDTO{
		ID:      reader.ID,
		Name:    reader.Name,
		Surname: reader.Surname,
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary Delete all readers
// @Tags readers
// @Success 204
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readers/ [delete]
func (h *ReadersHandler) DeleteAll(c *gin.Context) {
	if !h.config.EnableDeleteReaders {
		c.JSON(http.StatusForbidden, gin.H{"error": "DELETE /readers endpoint is disabled"})
		return
	}

	if err := h.repo.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete readers"})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get reader by ID
// @Tags readers
// @Produce json
// @Param id path int true "Reader ID"
// @Success 200 {object} dto.ReaderResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readers/{id} [get]
func (h *ReadersHandler) GetByID(c *gin.Context) {
	if !h.config.EnableGetReaders {
		c.JSON(http.StatusForbidden, gin.H{"error": "GET /readers endpoint is disabled"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	reader, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reader"})
		}
		return
	}

	books := make([]dto.BookResponseDTO, len(reader.CurrentlyReading))
	for i, book := range reader.CurrentlyReading {
		books[i] = dto.BookResponseDTO{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			UserID:      book.UserID,
			Username:    book.User.Username,
		}
	}

	response := dto.ReaderResponseDTO{
		ID:               reader.ID,
		Name:             reader.Name,
		Surname:          reader.Surname,
		CurrentlyReading: books,
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Update reader by ID
// @Tags readers
// @Accept json
// @Param id path int true "Reader ID"
// @Param reader body dto.ReaderUpdateDTO true "Updated reader data"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readers/{id} [put]
func (h *ReadersHandler) Update(c *gin.Context) {
	if !h.config.EnablePutReaders {
		c.JSON(http.StatusForbidden, gin.H{"error": "PUT /readers endpoint is disabled"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	reader, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reader"})
		}
		return
	}

	var readerDTO dto.ReaderUpdateDTO
	if err := c.ShouldBindJSON(&readerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if err := h.validator.ValidateStruct(readerDTO); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidationErrors(err))
		return
	}

	reader.Name = readerDTO.Name
	reader.Surname = readerDTO.Surname

	if err := h.repo.Update(reader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reader"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Delete reader by ID
// @Tags readers
// @Param id path int true "Reader ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readers/{id} [delete]
func (h *ReadersHandler) Delete(c *gin.Context) {
	if !h.config.EnableDeleteReaders {
		c.JSON(http.StatusForbidden, gin.H{"error": "DELETE /readers endpoint is disabled"})
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reader"})
		}
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reader"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddCurrentlyReading adds a book to reader's currently reading list
func (h *ReadersHandler) AddCurrentlyReading(c *gin.Context) {
	readerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book := &models.Book{}
	book.ID = uint(bookID)

	if err := h.repo.AddCurrentlyReading(uint(readerID), book); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to reading list"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveCurrentlyReading removes a book from reader's currently reading list
func (h *ReadersHandler) RemoveCurrentlyReading(c *gin.Context) {
	readerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.repo.RemoveCurrentlyReading(uint(readerID), uint(bookID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove book from reading list"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
