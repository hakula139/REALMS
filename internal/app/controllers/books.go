package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hakula139/REALMS/internal/app/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// ErrBookNotFound occurs when the queried book is not found
var ErrBookNotFound = errors.New("database: book not found")

// AddBookInput is a schema that validates input to prevent invalid requests
// ID will be generated automatically
type AddBookInput struct {
	Title     string `json:"title" binding:"required"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"isbn"`
}

// UpdateBookInput is a schema that validates input to prevent invalid requests
type UpdateBookInput struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"isbn"`
}

// RemoveBookInput is a schema that validates input to prevent invalid requests
type RemoveBookInput struct {
	Message string `json:"message"`
}

// FindBookInput is a schema that validates input to prevent invalid requests
type FindBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

// AddBook adds a new book to the library
// POST /admin/books
func AddBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validates input
	var input AddBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		Title:     input.Title,
		Author:    input.Author,
		Publisher: input.Publisher,
		ISBN:      input.ISBN,
	}
	db.Create(&book)

	logger := c.MustGet("logger").(*zap.SugaredLogger)
	logger.Infof("Added book %v", book.ID)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// UpdateBook updates data of a book
// PATCH /admin/books/:id
func UpdateBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var book models.Book
	bookID := c.Param("id")
	if err := db.Where("id = ?", bookID).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookNotFound.Error()})
		return
	}

	// Validates input
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&book).Updates(input)

	logger := c.MustGet("logger").(*zap.SugaredLogger)
	logger.Infof("Updated book %v", bookID)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// RemoveBook removes a book from the library
// DELETE /admin/books/:id
func RemoveBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var book models.Book
	bookID := c.Param("id")
	if err := db.Where("id = ?", bookID).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookNotFound.Error()})
		return
	}

	// Validates input
	var input RemoveBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Delete(&book)

	logger := c.MustGet("logger").(*zap.SugaredLogger)
	if input.Message == "" {
		logger.Infof("Removed book %v", bookID)
	} else {
		logger.Infof("Removed book %v with explanation: %v", bookID, input.Message)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// ShowBooks shows all books in the library
// GET /books
func ShowBooks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var books []models.Book
	db.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// ShowBook shows the book of given ID
// GET /books/:id
func ShowBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var book models.Book
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// FindBooks finds books by title / author / ISBN
// POST /books/find
func FindBooks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validates input
	var input FindBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var books []models.Book
	chain := db
	if ISBN := input.ISBN; ISBN != "" {
		chain = chain.Where("ISBN = ?", ISBN)
	}
	if author := input.Author; author != "" {
		chain = chain.Where("author = ?", author)
	}
	if title := input.Title; title != "" {
		chain = chain.Where("title LIKE ?", "%"+title+"%")
	}
	chain.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}
