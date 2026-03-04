package handler

import (
	"net/http"

	"github.com/airlangga-hub/library/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
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
