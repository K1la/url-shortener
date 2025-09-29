package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/K1la/url-shortener/internal/model"
	"github.com/go-redis/redis/v8"
)

func (s *Service) GetShortURL(ctx context.Context, redirectUrl model.RedirectClicks) (*model.URL, error) {
	url, err := s.cache.Get(redirectUrl.ShortURL)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get short url from cache: %w", err)
	}
	if errors.Is(err, redis.Nil) {
		return s.repo.GetShortURL(ctx, redirectUrl.ShortURL)
	}

	var URL model.URL
	URL.ShortURL = redirectUrl.ShortURL
	URL.URL = url

	return &URL, nil
}
