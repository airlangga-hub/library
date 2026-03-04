package repository

import (
	"fmt"

	"github.com/airlangga-hub/library/service"
	"gorm.io/gorm"
)

type repository struct {
	DB              *gorm.DB
	MailjetURL      string
	MailjetUsername string
	MailjetPassword string
	MailjetSender   string
}

func NewRepository(db *gorm.DB, mailjetURL, mailjetUsername, mailjetPassword, mailjetSender string) *repository {
	return &repository{
		DB:              db,
		MailjetURL:      mailjetURL,
		MailjetUsername: mailjetUsername,
		MailjetPassword: mailjetPassword,
		MailjetSender:   mailjetSender,
	}
}

func (r *repository) CreateUser(user service.User) (service.User, error) {
	u := User{
		FullName: user.FullName,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := r.DB.Create(&u).Error; err != nil {
		return service.User{}, fmt.Errorf("repo.CreateUser: %w", err)
	}

	return user, nil
}
