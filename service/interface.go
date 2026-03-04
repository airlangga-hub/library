package service

type Repository interface {
	SendEmail(to, subject, textPart string) error
	CreateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
	GetRents(userID int) ([]Rent, error)
}

type Service interface {
	Register(user User) (User, error)
	Login(email, password string) (string, error)
	GetRents(userID int) ([]Rent, error)
}