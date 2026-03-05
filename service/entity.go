package service

import "time"

type User struct {
	ID        int    `json:"-"`
	Admin     bool   `json:"-"`
	FullName  string `json:"full_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"-"`
	TotalRent int    `json:"total_rent,omitempty"`
	TotalBook int    `json:"total_book,omitempty"`
}

type Category struct {
	Name  string `json:"name,omitempty"`
	Books []Book `json:"books,omitempty"`
}

type Book struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author_id,omitempty"`
	Category    string `json:"category_id,omitempty"`
}

type Rent struct {
	BookTitle       string     `json:"book_title,omitempty"`
	BookDescription string     `json:"book_description,omitempty"`
	BookAuthor      string     `json:"book_author,omitempty"`
	BookCategory    string     `json:"book_category,omitempty"`
	RentDate        time.Time  `json:"rent_date"`
	DueDate         time.Time  `json:"due_date"`
	ReturnDate      *time.Time `json:"return_date,omitempty"`
	Fine            int        `json:"fine,omitempty"`
}
