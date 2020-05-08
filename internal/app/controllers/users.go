package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hakula139/REALMS/internal/app/models"
	"github.com/jinzhu/gorm"
)

// ErrUserNotFound occurs when the user is not found
var ErrUserNotFound = errors.New("database: user not found")

// AddUserInput is a schema that validates input to prevent invalid requests
// ID will be generated automatically
type AddUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Level    uint   `json:"level" binding:"required"`
}

// UpdateUserInput is a schema that validates input to prevent invalid requests
type UpdateUserInput struct {
	Password string `json:"password"`
	Level    uint   `json:"level"`
}

// AddUser adds a new user to the database
// POST /admin/users
func AddUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validates input
	var input AddUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: input.Password,
		Level:    input.Level}
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser updates data of a user
// PATCH /admin/users/:id
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrUserNotFound.Error()})
		return
	}

	// Validates input
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.EncryptPassword(&input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// RemoveUser removes a user from the database
// DELETE /admin/users/:id
func RemoveUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrUserNotFound.Error()})
		return
	}

	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// ShowUsers shows all users in the library
// GET /admin/users
func ShowUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []models.User
	db.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// ShowUser shows the user of given ID
// GET /admin/users/:id
func ShowUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrUserNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
