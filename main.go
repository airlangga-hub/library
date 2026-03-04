package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airlangga-hub/library/handler"
	"github.com/airlangga-hub/library/helper"
	"github.com/airlangga-hub/library/repository"
	"github.com/airlangga-hub/library/service"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

	if err := db.Exec("SET search_path TO library").Error; err != nil {
		log.Println("set search path failed:", err)
		return
	}
	
	if err := db.AutoMigrate(
		&repository.User{},
		&repository.Category{},
		&repository.Book{},
		&repository.Rent{},
	); err != nil {
		log.Println("db automigrate failed:", err)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	repo := repository.NewRepository(db, apiKey)
	svc := service.NewService(repo, []byte(jwtSecret))
	h := handler.NewHandler(svc, validate)

	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(helper.MyClaims)
		},
		SigningKey: []byte(jwtSecret),
	}

	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogMethod:   true,
		LogURI:      true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				errToLog := v.Error
				if internal := errors.Unwrap(v.Error); internal != nil {
					errToLog = internal
				}
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("error", errToLog.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// users
	users := e.Group("/users")
	// users public endpoints
	users.POST("/register", h.Register)
	users.POST("/login", h.Login)
	// users private endpoints
	users.POST("/topup", h.TopUp, echojwt.WithConfig(config))
	users.GET("/rents", h.GetRents)

	// books
	books := e.Group("/books", echojwt.WithConfig(config))
	books.POST("/rent", h.RentBook)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":" + port,
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
