package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/K1la/url-shortener/internal/model"
)

var (
	ErrShortURLNotFound = errors.New("short url not found")
)

func (r *Repository) GetShortURL(ctx context.Context, shortUrl string) (*model.URL, error) {
	query := `
		SELECT id, url, short_url
		FROM url
		WHERE short_url = $1
	`

	var url model.URL
	err := r.db.Master.QueryRowContext(ctx, query, shortUrl).
		Scan(&url.ID, &url.ShortURL, &url.ShortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrShortURLNotFound
		}

		return nil, fmt.Errorf("get short url: %w", err)
	}
	return &url, nil
}
