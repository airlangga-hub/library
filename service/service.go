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

	textPart := fmt.Sprintf("Hi %s!\n\nThanks for registering to Library FTGO 14!\nWe hope you're doing well!\n\nBest regards,\nLibraryFTGO 14", user.FullName)

	go func() {
		if err := s.Repo.SendEmail(user.Email, "Register Success", textPart); err != nil {
			slog.Error("Send Email Failed!!!", slog.Any("error", err))
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

	token, err := helper.MakeJWT(user.ID, email, s.JWTSecret)
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

func (s *service) RentBook(userID, bookID, duration int) (Rent, error) {
	createdAt := time.Now()
	returnDate := createdAt.Add(24 * time.Hour * time.Duration(duration))

	rent, err := s.Repo.CreateRent(userID, bookID, createdAt, returnDate)
	if err != nil {
		return Rent{}, fmt.Errorf("service.RentBook: %w", err)
	}
	return rent, nil
}
