package repository

import (
	"context"
	"fmt"
	"github.com/K1la/url-shortener/internal/model"
)

func (r *Repository) CreateShortURL(ctx context.Context, url model.URL) (*model.URL, error) {
	query := `
		INSERT INTO urls (url, short_url)
		VALUES ($1, $2)
		RETURNING id, short_url, created_at;
    `

	err := r.db.Master.QueryRowContext(
		ctx, query, url.URL, url.ShortURL,
	).Scan(&url.ID, &url.ShortURL, &url.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert url into db: %w", err)
	}

	return &url, nil

}
