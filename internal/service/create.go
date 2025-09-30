package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/K1la/url-shortener/internal/model"
	"github.com/K1la/url-shortener/internal/repository"
	"github.com/wb-go/wbf/zlog"
)

var (
	ErrShortURLAlreadyExists = errors.New("short url already exists")
)

func (s *Service) CreateShortURL(ctx context.Context, url model.URL) (*model.URL, error) {
	url.URL = validateURL(url.URL)

	// If no shortURl provided
	if url.ShortURL == "" {
		for {
			url.ShortURL = generateShortLink()
			_, err := s.repo.GetShortURL(ctx, url.ShortURL)
			if errors.Is(err, repository.ErrShortURLNotFound) {
				break // unique shortURL found
			}
		}
	} else {
		// If shortURL provided check if it exists
		_, err := s.repo.GetShortURL(ctx, url.ShortURL)
		if err != nil && !errors.Is(err, repository.ErrShortURLNotFound) {
			return nil, fmt.Errorf("failed to check existing provided shortURL: %w", err)
		}
		if err == nil {
			return nil, ErrShortURLAlreadyExists
		}

	}

	res, err := s.repo.CreateShortURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to create short url: %w", err)
	}

	//b, err := json.Marshal(res)
	//if err != nil {
	//	return nil, fmt.Errorf("marshal short url: %w", err)
	//}
	//
	//err = s.cache.Set(url.ShortURL, string(b))

	if b, err := json.Marshal(res); err == nil {
		if err = s.cache.Set(url.ShortURL, string(b)); err != nil {
			zlog.Logger.Error().Err(err).Str("url", url.ShortURL).Msg("failed to cache short url")
		}
	}

	return res, err
}

func (s *Service) SaveAnalytics(ctx context.Context, rURL *model.RedirectClicks) (string, error) {
	id, err := s.repo.SaveAnalytics(ctx, rURL)
	if err != nil {
		return "", err
	}

	rURL.ID = id

	go func() {
		bgCtx := context.Background()

		summary, err := s.GetAnalyticsSummary(bgCtx, rURL.ShortURL)
		if err != nil {
			zlog.Logger.Error().Err(err).Str("shorturl", rURL.ShortURL).Msg("failed to get analytics summary from cache")
			return
		}

		if b, err := json.Marshal(summary); err == nil {
			if err := s.cache.Set("analytics:"+rURL.ShortURL, string(b)); err != nil {
				zlog.Logger.Error().Err(err).Str("shorturl", rURL.ShortURL).Msg("failed to cache aggregated analytics")
			}
		}

		zlog.Logger.Info().Str("shorturl", rURL.ShortURL).Msg("saved analytics")
	}()

	return rURL.ID, nil
}
