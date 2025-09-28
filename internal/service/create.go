package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/K1la/url-shortener/internal/model"
	"github.com/K1la/url-shortener/internal/repository"
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

	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshal short url: %w", err)
	}

	err = s.cache.Set(url.URL, string(b))

	return res, err

	//shortUrl, err := s.cache.Get(url.URL)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to fetch short url from cache: %w", err)
	//}
	//
	//if errors.Is(err, redis.Nil) {
	//	attempts := 3
	//	for range attempts {
	//		url.ShortURL = generateShortLink()
	//		urlInfo, err := s.repo.CreateShortURL(ctx, url)
	//		if err == nil {
	//			return urlInfo, nil
	//		}
	//	}
	//	return nil, err
	//}
	//
	//var urlInfo model.URL
	//urlInfo.URL = url.URL
	//urlInfo.ShortURL = shortUrl
	//
	//return &urlInfo, nil
}
