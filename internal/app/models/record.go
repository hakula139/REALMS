package models

import (
	"time"
)

// Record is stored when a user borrows a book from the library,
// and is soft deleted when the book is returned
type Record struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id" gorm:"NOT NULL"`
	BookID      uint       `json:"book_id" gorm:"NOT NULL"`
	ReturnDate  time.Time  `json:"return_date" gorm:"NOT NULL"`
	ExtendTimes uint       `json:"extend_times" gorm:"NOT NULL"`
	DeletedAt   *time.Time `json:"real_return_date"`
}
