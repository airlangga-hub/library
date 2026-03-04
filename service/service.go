package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	Repo      Repository
	JWTSecret []byte
}

func NewService(repo Repository, jwtSecret []byte) *service {
	return &service{Repo: repo, JWTSecret: jwtSecret}
}

func (s *service) Register(user User) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, fmt.Errorf("service.Register: %w", err)
	}
	
	user.Password = string(hashedPassword)
	
	user, err = s.Repo.Register(user)
	if err != nil {
		return User{}, fmt.Errorf("service.Register: %w", err)
	}
	
	return user, nil
}