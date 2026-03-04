package service

type Repository interface {
	CreateUser(user User) (User, error)
}

type Service interface {
	Register(user User) (User, error)
}