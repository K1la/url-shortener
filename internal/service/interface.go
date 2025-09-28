package service

import (
	"context"
	"github.com/K1la/url-shortener/internal/model"
)

type RepositoryI interface {
	CreateShortURL(context.Context, model.URL) (*model.URL, error)
	GetShortURL(context.Context, string) (*model.URL, error)
}

type CacheI interface {
	Get(string) (string, error)
	Set(string, interface{}) error
}
