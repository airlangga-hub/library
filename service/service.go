package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/airlangga-hub/library/helper"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	Repo      Repository
	JWTSecret []byte
	Logger    *slog.Logger
}

func NewService(repo Repository, jwtSecret []byte, logger *slog.Logger) *service {
	return &service{
		Repo:      repo,
		JWTSecret: jwtSecret,
		Logger:    logger,
	}
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

	go func() {
		textPart := fmt.Sprintf("Hi %s!\n\nThanks for registering to Library FTGO 14!\nWe hope you're doing well!\n\nBest regards,\nLibrary FTGO 14", user.FullName)

		if err := s.Repo.SendEmail(user.Email, "Register Success", textPart); err != nil {
			s.Logger.Error("Send Email Failed!!!", slog.Any("error", err))
		}
	}()

	return user, nil
}

func (s *service) Login(email, password string) (string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("service.Login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("service.Login: %w", err)
	}

	var token string

	if user.Admin {
		token, err = helper.MakeJWT(true, user.ID, email, s.JWTSecret)
	} else {
		token, err = helper.MakeJWT(false, user.ID, email, s.JWTSecret)
	}

	if err != nil {
		return "", fmt.Errorf("service.Login: %w", err)
	}

	return token, nil
}

func (s *service) GetRents(userID int) ([]Rent, error) {
	rents, err := s.Repo.GetRents(userID)
	if err != nil {
		return nil, fmt.Errorf("service.GetRents: %w", err)
	}
	return rents, nil
}

func (s *service) RentBook(email string, userID, bookID, duration int) (Rent, error) {
	createdAt := time.Now()
	dueDate := createdAt.Add(24 * time.Hour * time.Duration(duration))

	rent, err := s.Repo.CreateRent(userID, bookID, createdAt, dueDate)
	if err != nil {
		return Rent{}, fmt.Errorf("service.RentBook: %w", err)
	}

	go func() {
		textPart := "Thanks for choosing Library FTGO 14 for your recent book rental.\n\nWe hope our books provide you with the resources you were looking for.\n\nBest regards,\nLibrary FTGO 14"

		if err := s.Repo.SendEmail(email, "Thank You for Choosing Library FTGO 14", textPart); err != nil {
			s.Logger.Error("Send Email Failed!!!", slog.Any("error", err))
		}
	}()

	return rent, nil
}

func (s *service) GetBooks() ([]Book, error) {
	books, err := s.Repo.GetBooks()
	if err != nil {
		return nil, fmt.Errorf("service.GetBooks: %w", err)
	}
	return books, nil
}

func (s *service) AdminGetRentsReport() ([]User, error) {
	users, err := s.Repo.AdminGetRentsReport()
	if err != nil {
		return nil, fmt.Errorf("service.AdminGetRentsReport: %w", err)
	}
	return users, nil
}

func (s *service) AdminGetAuthorsReport() ([]User, error) {
	users, err := s.Repo.AdminGetAuthorsReport()
	if err != nil {
		return nil, fmt.Errorf("service.AdminGetAuthorsReport: %w", err)
	}
	return users, nil
}

func (s *service) ReturnBook(userID, bookID int) (Rent, error) {
	rent, err := s.Repo.ReturnBook(userID, bookID)
	if err != nil {
		return Rent{}, fmt.Errorf("service.ReturnBook: %w", err)
	}
	return rent, nil
}
