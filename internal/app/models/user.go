package models

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ErrUsernameRequired occurs when the username field is left blank
var ErrUsernameRequired = errors.New("validate: username required")

// ErrPasswordRequired occurs when the password field is left blank
var ErrPasswordRequired = errors.New("validate: password required")

// User is a person who has access to the library
// Level indicates the user's privilege
// 0: Banned
// 1: User
// 2: Admin
// 3: Super Admin
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username" gorm:"NOT NULL; UNIQUE"`
	Password string `json:"password" gorm:"NOT NULL"`
	Level    uint   `json:"level" gorm:"NOT NULL"`
}

func hash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

// TrimUsername removes whitespaces in the username
func TrimUsername(username *string) {
	trimmed := strings.Trim(*username, " ")
	*username = trimmed
}

// EncryptPassword encrypts the password
func EncryptPassword(pass *string) error {
	hash, err := hash(*pass)
	if err != nil {
		return err
	}
	*pass = string(hash)
	return nil
}

// Validate checks if the user has input the required fields
func Validate(username, password string) error {
	TrimUsername(&username)
	if username == "" {
		return ErrUsernameRequired
	}
	if password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// VerifyPassword verifies that the password matches the hash
func VerifyPassword(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

// IsAdmin checks if the user has admin privilege
func (u *User) IsAdmin() bool {
	return u.Level >= 2
}

// BeforeSave trims the username and encrypts the password before saving user data
func (u *User) BeforeSave() error {
	TrimUsername(&u.Username)
	return EncryptPassword(&u.Password)
}

// Validate checks if the user has input the required fields
func (u *User) Validate() error {
	return Validate(u.Username, u.Password)
}
