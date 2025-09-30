package handler

import (
	"context"
	"github.com/K1la/url-shortener/internal/model"
)

type ServiceI interface {
	CreateShortURL(context.Context, model.URL) (*model.URL, error)
	GetShortURL(context.Context, model.RedirectClicks) (*model.URL, error)
	SaveAnalytics(context.Context, *model.RedirectClicks) (string, error)
	GetAnalyticsSummary(context.Context, string) (*model.SummaryOfAnalytics, error)
	InvalidateAnalyticsCache(context.Context, string) error
}
