package handler

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RentBookRequest struct {
	BookID   int `json:"book_id" validate:"required,gt=0"`
	Duration int `json:"duration" validate:"required,gt=0,lt=15"`
}
