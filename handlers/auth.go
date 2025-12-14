package handlers

import (
	"errors"
	"lab1/dto"
	"lab1/middleware"
	"lab1/models"
	"lab1/repository"
	"lab1/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo  repository.UserRepository
	validator *validation.Validator
}

func NewAuthHandler(userRepo repository.UserRepository, validator *validation.Validator) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		validator: validator,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidationErrors(err))
		return
	}

	// Check if username already exists
	existingUser, err := h.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	existingUser, err = h.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "user",
	}

	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := dto.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidationErrors(err))
		return
	}

	// Get user by username
	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := dto.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userRepo.GetByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusOK, response)
}
