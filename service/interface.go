package service

import "time"

type Repository interface {
	SendEmail(to, subject, textPart string) error
	CreateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
	GetRents(userID int) ([]Rent, error)
	CreateRent(userID, bookID int, createdAt, returnDate time.Time) (Rent, error)
	GetBooks() ([]Book, error)
	AdminGetRentsReport() ([]UserRentReport, error)
	AdminGetAuthorsReport() ([]AuthorBookReport, error)
	ReturnBook(userID, bookID int) (Rent, error)
}

type Service interface {
	Register(user User) (User, error)
	Login(email, password string) (string, error)
	GetRents(userID int) ([]Rent, error)
	RentBook(email string, userID, bookID, duration int) (Rent, error)
	GetBooks() ([]Book, error)
	AdminGetRentsReport() ([]UserRentReport, error)
	AdminGetAuthorsReport() ([]AuthorBookReport, error)
	ReturnBook(userID, bookID int) (Rent, error)
}
