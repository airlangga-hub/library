package main

import (
	"log"
	"os"

	"github.com/airlangga-hub/library/handler"
	"github.com/airlangga-hub/library/repository"
	"github.com/airlangga-hub/library/service"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	godotenv.Load()
	
	dsn := os.Getenv("DSN")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("JWT_SECRET")
	if dsn == "" || port == "" || jwtSecret == "" || apiKey == "" {
		log.Fatalln("env variable missing.")
	}
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalln("db open failed:", err)
	}
	
	validate := validator.New(validator.WithRequiredStructEnabled())
	
	repo := repository.NewRepository(db, apiKey)
	svc := service.NewService(repo, []byte(jwtSecret))
	h := handler.NewHandler(svc, validate)
}