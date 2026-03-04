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

	user, err = s.Repo.CreateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("service.Register: %w", err)
	}

	textPart := fmt.Sprintf("Hi %s!\n\nThanks for registering to Library FTGO 14!\nWe hope you're doing well!\n\nBest regards,\nLibraryFTGO 14", user.FullName)

	go s.Repo.SendEmail(user.Email, "Register Success", textPart)

	return user, nil
}
