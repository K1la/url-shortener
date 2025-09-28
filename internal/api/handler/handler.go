package handler

import (
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service ServiceI
	valid   *validator.Validate
}

func New(s ServiceI, v *validator.Validate) *Handler {
	return &Handler{service: s, valid: v}
}

type CreateRequest struct {
	URL      string `json:"url"       validate:"required"`
	ShortURL string `json:"short_url" validate:"required"`
}
