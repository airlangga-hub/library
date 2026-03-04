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

func (r *repository) GetUserByEmail(email string) (service.User, error) {
	user := User{}
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return service.User{}, fmt.Errorf("repo.GetUserByEmail: %w", err)
	}
	return service.User{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *repository) GetRents(userID int) ([]service.Rent, error) {
	rents := make([]Rent, 0, 16)

	res := r.DB.
		Where("user_id = ?", userID).
		Joins("Book").
		Joins("Book.Category").
		Joins("Book.Author").
		Find(&rents)

	if err := res.Error; err != nil {
		return nil, fmt.Errorf("repo.GetRents: %w", err)
	}
	if len(rents) == 0 {
		return nil, fmt.Errorf("repo.GetRents: %w", gorm.ErrRecordNotFound)
	}

	rrents := make([]service.Rent, len(rents))
	for i, r := range rents {
		rrents[i] = service.Rent{
			BookTitle:       r.Book.Title,
			BookDescription: r.Book.Description,
			BookAuthor:      r.Book.Author.FullName,
			BookCategory:    r.Book.Category.Name,
			RentDate:        r.RentDate,
		}
	}
	
	return rrents, nil
}
