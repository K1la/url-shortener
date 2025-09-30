package service

import (
	"context"
	"github.com/K1la/url-shortener/internal/model"
)

type RepositoryI interface {
	CreateShortURL(context.Context, model.URL) (*model.URL, error)
	GetShortURL(context.Context, string) (*model.URL, error)
	SaveAnalytics(context.Context, *model.RedirectClicks) (string, error)
	CountClicks(context.Context, string) (int, error)
	CountClicksByDay(context.Context, string) (map[string]int, error)
	CountClicksByMonth(context.Context, string) (map[string]int, error)
	CountClicksByUserAgent(context.Context, string) (map[string]int, error)
}

type CacheI interface {
	Get(string) (string, error)
	Set(string, interface{}) error
	Delete(string) error
}
