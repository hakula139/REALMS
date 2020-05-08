package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hakula139/REALMS/internal/app/models"
	"github.com/jinzhu/gorm"
)

const userkey = "user"

// ErrUnauthorized occurs when the user is unauthorized to perform the operation
var ErrUnauthorized = errors.New("auth: unauthorized")

// ErrUserNotExist occurs when the username is not found
var ErrUserNotExist = errors.New("auth: user not exist")

// ErrAuthFailed occurs when the password is incorrect
var ErrAuthFailed = errors.New("auth: incorrect password")

// ErrAlreadyLoggedIn occurs when the user has already logged in
var ErrAlreadyLoggedIn = errors.New("auth: already logged in")

// ErrSaveSessionFailed occurs when failed to save session
var ErrSaveSessionFailed = errors.New("auth: failed to save session")

// ErrInvalidSession occurs when the session token is not found or invalid
var ErrInvalidSession = errors.New("auth: invalid session token")

// AuthRequired is a middleware that validates the session
// User privilege required
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get(userkey)
	if uid == nil {
		// Aborts the request
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}
	c.Next()
}

// AdminRequired is a middleware that validates the session
// Admin privilege required
func AdminRequired(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get(userkey)
	if uid == nil {
		// Aborts the request
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	// Checks if the user has admin privilege
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", uid).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUserNotExist.Error()})
		return
	}
	if !user.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	c.Next()
}

// Login verifies user identity and saves the session token
// POST /login
func Login(c *gin.Context) {
	session := sessions.Default(c)
	if uid := session.Get(userkey); uid != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrAlreadyLoggedIn.Error()})
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validates input
	if err := models.Validate(username, password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifies username and password
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrUserNotExist.Error()})
		return
	}
	if err := models.VerifyPassword(user.Password, password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrAuthFailed.Error()})
		return
	}

	// Saves the user ID in the session
	session.Set(userkey, user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrSaveSessionFailed.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// Logout removes the session token
// GET /logout
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	if uid := session.Get(userkey); uid == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidSession.Error()})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrSaveSessionFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// Me shows the currently logged-in user
// GET /user/me
func Me(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get(userkey)
	c.JSON(http.StatusOK, gin.H{"user": uid})
}

// Status shows the current log-in status
// GET /user/status
func Status(c *gin.Context) {
	session := sessions.Default(c)
	if uid := session.Get(userkey); uid == nil {
		c.JSON(http.StatusOK, gin.H{"status": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true})
}
