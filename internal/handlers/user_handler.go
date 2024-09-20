package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"write-stream-go/internal/models"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type CreateUserInput struct {
	SocialID  string  `json:"SocialID" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Name      string  `json:"name" binding:"required"`
	AvatarURL *string `json:"avatarUrl"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		ID:        uuid.New(),
		SocialID:  input.SocialID,
		Email:     input.Email,
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if input.AvatarURL != nil {
		user.AvatarURL = *input.AvatarURL
	}

	if err := h.DB.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
