package models

// Book is the basic object stored in the database
type Book struct {
	ID        uint   `json:"id"`
	Title     string `json:"title" gorm:"NOT NULL"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"ISBN"`
}
