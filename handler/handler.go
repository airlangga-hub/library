package handler

import (
	"github.com/airlangga-hub/library/service"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	Svc      service.Service
	Validate *validator.Validate
}

func NewHandler(svc service.Service, val *validator.Validate) *handler {
	return &handler{Svc: svc, Validate: val}
}