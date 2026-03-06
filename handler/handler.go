package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/airlangga-hub/library/helper"
	"github.com/airlangga-hub/library/service"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type handler struct {
	Svc      service.Service
	Validate *validator.Validate
}

func NewHandler(svc service.Service, val *validator.Validate) *handler {
	return &handler{Svc: svc, Validate: val}
}

func (h *handler) Register(c *echo.Context) error {
	var payload RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	u := service.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
	}

	user, err := h.Svc.Register(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "register failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    user,
	})
}

func (h *handler) Login(c *echo.Context) error {
	var payload LoginRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	token, err := h.Svc.Login(payload.Email, payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid credentials").Wrap(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *handler) GetRents(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	rents, err := h.Svc.GetRents(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no rents found for this user").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "get rents failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    rents,
	})
}

func (h *handler) RentBook(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload RentBookRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				if e.Field() == "Duration" && e.Tag() == "lt" {
					return echo.NewHTTPError(http.StatusBadRequest, "rent duration maximum 14 days").Wrap(e)
				}
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	rent, err := h.Svc.RentBook(claims.Subject, claims.UserID, payload.BookID, payload.Duration)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "book unavailable for now").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "rent book failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    rent,
	})
}

func (h *handler) GetBooks(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	_, ok = token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	books, err := h.Svc.GetBooks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "get books failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    books,
	})
}

func (h *handler) AdminGetRentsReport(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok || !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	userRents, err := h.Svc.AdminGetRentsReport()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "get rents report failed")
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    userRents,
	})
}

func (h *handler) AdminGetAuthorsReport(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok || !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	userBooks, err := h.Svc.AdminGetRentsReport()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "get books report failed")
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    userBooks,
	})
}

func (h *handler) ReturnBook(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok || !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "book id must be integer").Wrap(err)
	}

	rent, err := h.Svc.ReturnBook(claims.UserID, bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no rent with that book id found").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusOK, "return book failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    rent,
	})
}
