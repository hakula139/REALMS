package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hakula139/REALMS/internal/app/config"
	"github.com/hakula139/REALMS/internal/app/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

const day = time.Hour * 24

// ErrBookBorrowed occurs when the user wants to borrow a book which has been
// borrowed before
var ErrBookBorrowed = errors.New("library: book already borrowed")

// ErrBookNotBorrowed occurs when the user wants to return a book which has
// not been borrowed before
var ErrBookNotBorrowed = errors.New("library: book not borrowed")

// ErrRecordNotFound occurs when the record is not found
var ErrRecordNotFound = errors.New("database: record not found")

// ErrExceedMaxExtendTimes occurs when the user has extended the deadline too
// many times
var ErrExceedMaxExtendTimes = errors.New("library: extended too many times")

// AddRecordInput is a schema that validates input to prevent invalid requests
// ID, UserID, BookID, ExtendTimes will be generated automatically
// BorrowDate will be set to current date if left blank
type AddRecordInput struct {
	BorrowDate time.Time `json:"borrow_date"`
}

// BorrowBook adds a new record to the database
// POST /user/books/:id
func BorrowBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	libcfg := c.MustGet("libcfg").(config.LibraryConfig)

	// Gets user ID
	session := sessions.Default(c)
	userID := session.Get(userkey)

	// Checks if the user should be suspended
	var count uint
	today := time.Now().Local()
	db.Model(&models.Record{}).Where("user_id = ? AND return_date < ?", userID, today).Count(&count)
	if count >= libcfg.MaxOverdueBooks {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	// Gets book ID and checks if the book exists
	count = 0
	bookID := c.Param("id")
	db.Model(&models.Book{}).Where("id = ?", bookID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookNotFound.Error()})
		return
	}

	// Checks if the book has been borrowed before
	count = 0
	db.Model(&models.Record{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookBorrowed.Error()})
		return
	}

	// Validates input
	var input AddRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculates return date
	var returnDate time.Time
	expire := time.Duration(libcfg.BorrowExpireDays) * day
	if input.BorrowDate.IsZero() {
		returnDate = today.Add(expire)
	} else {
		returnDate = input.BorrowDate.Add(expire)
	}

	bookIDInt, _ := strconv.Atoi(bookID)
	bookIDUint := uint(bookIDInt)
	userIDUint, _ := userID.(uint)

	record := models.Record{
		UserID:      userIDUint,
		BookID:      bookIDUint,
		ReturnDate:  returnDate,
		ExtendTimes: 0,
	}
	db.Create(&record)

	logger := c.MustGet("logger").(*zap.SugaredLogger)
	logger.Infof("User %v borrowed book %v", userID, bookID)

	c.JSON(http.StatusOK, gin.H{"data": record})
}

// ExtendDeadline extends the deadline to return a book
// PATCH /user/books/:id
func ExtendDeadline(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	libcfg := c.MustGet("libcfg").(config.LibraryConfig)

	// Gets user ID
	session := sessions.Default(c)
	userID := session.Get(userkey)

	var record models.Record
	bookID := c.Param("id")
	if err := db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&record).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrRecordNotFound.Error()})
		return
	}

	// Checks if the user has extended the deadline too many times
	if record.ExtendTimes >= libcfg.MaxExtendTimes {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrExceedMaxExtendTimes.Error()})
		return
	}

	extend := time.Duration(libcfg.DdlExtendDays) * day
	db.Model(&record).Updates(models.Record{
		ReturnDate:  record.ReturnDate.Add(extend),
		ExtendTimes: record.ExtendTimes + 1,
	})

	c.JSON(http.StatusOK, gin.H{"data": record})
}

// ReturnBook soft deletes the related record from the database
// DELETE /user/books/:id
func ReturnBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Gets user ID
	session := sessions.Default(c)
	userID := session.Get(userkey)

	// Gets book ID and checks if the book has been borrowed before
	var record models.Record
	bookID := c.Param("id")
	if err := db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&record).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookNotBorrowed.Error()})
		return
	}

	db.Delete(&record)

	logger := c.MustGet("logger").(*zap.SugaredLogger)
	logger.Infof("User %v returned book %v", userID, bookID)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// ShowBookList shows all books that the user has borrowed
// GET /user/books
func ShowBookList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	session := sessions.Default(c)
	userID := session.Get(userkey)

	var records []models.Record
	db.Order("return_date").Where("user_id = ?", userID).Find(&records)

	c.JSON(http.StatusOK, gin.H{"data": records})
}

// ShowBorrowed shows a book that the user has borrowed
// Deadline is also returned in data.return_date
// GET /user/books/:id
func ShowBorrowed(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Gets user ID
	session := sessions.Default(c)
	userID := session.Get(userkey)

	// Gets book ID and checks if the book has been borrowed before
	var record models.Record
	bookID := c.Param("id")
	if err := db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&record).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBookNotBorrowed.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": record})
}

// ShowOverdueList shows all overdue books that the user has borrowed
// GET /user/overdue
func ShowOverdueList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	session := sessions.Default(c)
	userID := session.Get(userkey)

	var records []models.Record
	today := time.Now().Local()
	db.Order("return_date").Where("user_id = ? AND return_date < ?", userID, today).Find(&records)

	c.JSON(http.StatusOK, gin.H{"data": records})
}

// ShowHistory shows all records of the user
// GET /user/history
func ShowHistory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	session := sessions.Default(c)
	userID := session.Get(userkey)

	var records []models.Record
	db.Order("return_date").Unscoped().Where("user_id = ?", userID).Find(&records)

	c.JSON(http.StatusOK, gin.H{"data": records})
}
