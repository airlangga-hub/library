package service

type Repository interface {
	Register(user User) (User, error)
}

type Service interface {
	Register(user User) (User, error)
}