package handler

import (
	"errors"
	"net/http"

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
