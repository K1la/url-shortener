package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/K1la/url-shortener/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) GetShortURL(ctx context.Context, redirectUrl model.RedirectClicks) (*model.URL, error) {
	var URL *model.URL
	url, err := s.cache.Get(redirectUrl.ShortURL)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get short url from cache: %w", err)
	}
	if errors.Is(err, redis.Nil) {
		zlog.Logger.Info().Str("rUrl", redirectUrl.ShortURL).Msg("short url not found in cache")
		URL, err = s.repo.GetShortURL(ctx, redirectUrl.ShortURL)
		if err != nil {
			return nil, err
		}

		//b, err := json.Marshal(URL)
		//if err != nil {
		//	return nil, fmt.Errorf("marshal short url: %w", err)
		//}
		//
		//err = s.cache.Set(URL.ShortURL, string(b))

		if b, err := json.Marshal(URL); err == nil {
			if err = s.cache.Set(URL.ShortURL, string(b)); err != nil {
				zlog.Logger.Error().Err(err).Str("url", URL.ShortURL).Msg("failed to cache short url")
			}
		}

		return URL, nil
	}
	zlog.Logger.Info().Str("url", url).Msg("short url found in cache!!!")

	err = json.Unmarshal([]byte(url), &URL)
	if err != nil {
		return nil, fmt.Errorf("unmarshal short url: %w", err)
	}

	return URL, nil
}

func (s *Service) GetAnalyticsSummary(ctx context.Context, shortURL string) (*model.SummaryOfAnalytics, error) {
	if str, err := s.cache.Get("analytics:" + shortURL); err == nil {
		var sum model.SummaryOfAnalytics
		if err = json.Unmarshal([]byte(str), &sum); err == nil {
			zlog.Logger.Info().Str("shorturl", shortURL).Interface("summary", sum).Msg("analytics summary from cache")
			return &sum, nil // found in cache
		}
	}

	clicks, err := s.repo.CountClicks(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("count clicks: %w", err)
	}

	daily, err := s.repo.CountClicksByDay(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("count clicks by day: %w", err)
	}

	monthly, err := s.repo.CountClicksByMonth(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("count clicks by month: %w", err)
	}

	ua, err := s.repo.CountClicksByUserAgent(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("count clicks by month: %w", err)
	}

	sum := &model.SummaryOfAnalytics{
		ShortUrl:    shortURL,
		TotalClicks: clicks,
		Daily:       daily,
		Monthly:     monthly,
		UserAgent:   ua,
	}

	zlog.Logger.Info().Str("shorturl", shortURL).Interface("summary", sum).Msg("analytics summary from db")

	if b, err := json.Marshal(sum); err == nil {
		if err = s.cache.Set("analytics:"+shortURL, string(b)); err != nil {
			zlog.Logger.Error().Err(err).Str("short url", shortURL).Msg("failed to cache summary analytics")
		}
	}

	return sum, nil
}

func (s *Service) InvalidateAnalyticsCache(ctx context.Context, shortURL string) error {
	// Delete analytics cache for this short URL
	cacheKey := "analytics:" + shortURL
	if err := s.cache.Delete(cacheKey); err != nil {
		zlog.Logger.Error().Err(err).Str("shortUrl", shortURL).Str("cacheKey", cacheKey).Msg("failed to delete analytics cache")
		return fmt.Errorf("failed to delete analytics cache: %w", err)
	}

	zlog.Logger.Info().Str("shortUrl", shortURL).Str("cacheKey", cacheKey).Msg("analytics cache invalidated")
	return nil
}
