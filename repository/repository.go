package repository

import "gorm.io/gorm"

type repository struct {
	DB     *gorm.DB
	APIKey string
}

func NewRepository(db *gorm.DB, apiKey string) *repository {
	return &repository{DB: db, APIKey: apiKey}
}