package handler

import (
	"context"
	"github.com/K1la/url-shortener/internal/model"
)

type ServiceI interface {
	CreateShortURL(context.Context, model.URL) (*model.URL, error)
}
