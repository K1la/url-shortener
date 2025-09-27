package handlers

import (
	"time"

	"github.com/K1la/url-shortener/internal/model"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service ServiceLinkI
	valid   *validator.Validate
}

func New(s ServiceI, v *validator.Validate) *Handler {
	return &Handler{service: s, valid: v}
}

type CreateRequest struct {
	`json:"" validate:"required"`
	`json:"" validate:"required"`
	`json:"" validate:"required"`
	`json:"" validate:"required"`
	`json:"" validate:"required"`
}
