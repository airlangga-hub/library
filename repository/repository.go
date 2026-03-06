package repository

import (
	"fmt"
	"time"

	"github.com/airlangga-hub/library/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Admin:    user.Admin,
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
		Omit("Book.Author.total_rent", "Book.Author.total_book").
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
			RentDate:        r.CreatedAt,
			DueDate:         r.DueDate,
			ReturnDate:      r.ReturnDate,
			Fine:            r.Fine,
			Active:          r.Active,
		}
	}

	return rrents, nil
}

func (r *repository) CreateRent(userID, bookID int, createdAt, dueDate time.Time) (service.Rent, error) {
	rent := Rent{
		BookID:    bookID,
		UserID:    userID,
		CreatedAt: createdAt,
		DueDate:   dueDate,
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&Book{}).
			Where("id = ? AND available = true", bookID).
			Update("available", false)

		if err := res.Error; err != nil {
			return err
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		err := tx.Create(&rent).Error
		if err != nil {
			return err
		}

		err = r.DB.
			Joins("Book").
			Joins("Book.Author").
			Joins("Book.Category").
			Omit("Book.Author.total_rent", "Book.Author.total_book").
			First(&rent).
			Error

		return err
	})

	if err != nil {
		return service.Rent{}, fmt.Errorf("repo.CreateRent: %w", err)
	}

	return service.Rent{
		BookTitle:       rent.Book.Title,
		BookDescription: rent.Book.Description,
		BookAuthor:      rent.Book.Author.FullName,
		BookCategory:    rent.Book.Category.Name,
		RentDate:        rent.CreatedAt,
		DueDate:         rent.DueDate,
		Active:          rent.Active,
	}, nil
}

func (r *repository) ReturnBook(userID, bookID int) (service.Rent, error) {
	book := Book{ID: bookID}
	var rent Rent

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Joins("Book").
			Joins("Book.Author").
			Joins("Book.Category").
			Omit("Book.Author.total_rent", "Book.Author.total_book").
			Where("user_id = ? AND book_id = ? AND return_date IS NULL", userID, bookID).
			First(&rent).
			Error
		if err != nil {
			return err
		}

		fine := 0
		now := time.Now()
		hoursLate := int(now.Sub(rent.DueDate).Hours())

		if hoursLate >= 1 {
			fine = 2000 * hoursLate
		}

		rent.Fine = fine
		rent.ReturnDate = &now

		res := tx.Model(&rent).
			Updates(map[string]any{"fine": fine, "return_date": now, "active": false})
		if err := res.Error; err != nil {
			return err
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		res = tx.Model(&book).
			Where("available = false").
			Update("available", true)
		if err := res.Error; err != nil {
			return err
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	if err != nil {
		return service.Rent{}, fmt.Errorf("repo.ReturnBook: %w", err)
	}

	return service.Rent{
		BookTitle:       rent.Book.Title,
		BookDescription: rent.Book.Description,
		BookAuthor:      rent.Book.Author.FullName,
		BookCategory:    rent.Book.Category.Name,
		RentDate:        rent.CreatedAt,
		DueDate:         rent.DueDate,
		ReturnDate:      rent.ReturnDate,
		Fine:            rent.Fine,
		Active:          rent.Active,
	}, nil
}

func (r *repository) GetBooks() ([]service.Book, error) {
	books := make([]Book, 0, 16)

	err := r.DB.Joins("Author").Joins("Category").Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("repo.GetBooks: %w", err)
	}

	bbooks := make([]service.Book, len(books))
	for i, b := range books {
		bbooks[i] = service.Book{
			Title:       b.Title,
			Description: b.Description,
			Author:      b.Author.FullName,
			Category:    b.Category.Name,
		}
	}

	return bbooks, nil
}

func (r *repository) AdminGetRentsReport() ([]service.User, error) {
	users := make([]User, 0, 16)

	err := r.DB.
		Model(&User{}).
		Select("users.id, users.full_name, users.email, COUNT(rents.id) AS total_rent").
		Joins("LEFT JOIN rents ON rents.user_id = users.id").
		Group("users.id").
		Order("total_rent DESC").
		Find(&users).
		Error

	if err != nil {
		return nil, fmt.Errorf("repo.AdminGetRentsReport: %w", err)
	}

	uusers := make([]service.User, len(users))
	for i, u := range users {
		uusers[i] = service.User{
			FullName:  u.FullName,
			Email:     u.Email,
			TotalRent: u.TotalRent,
		}
	}

	return uusers, nil
}

func (r *repository) AdminGetAuthorsReport() ([]service.User, error) {
	users := make([]User, 0, 16)

	err := r.DB.
		Select("users.id, users.full_name, users.email, users.author, COUNT(books.id) AS total_book").
		Joins("LEFT JOIN books ON books.author_id = users.id").
		Where("users.author = true").
		Group("users.id").
		Find(&users).
		Error

	if err != nil {
		return nil, fmt.Errorf("repo.AdminGetAuthorsReport: %w", err)
	}

	uusers := make([]service.User, len(users))
	for i, u := range users {
		uusers[i] = service.User{
			FullName:  u.FullName,
			Email:     u.Email,
			TotalBook: u.TotalBook,
		}
	}

	return uusers, nil
}
