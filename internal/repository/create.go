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

func (r *Repository) SaveAnalytics(ctx context.Context, rUrl *model.RedirectClicks) (string, error) {
	query := `
	INSERT INTO analytics (
        short_url, user_agent, device_type, os, browser, ip_address
	) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`

	err := r.db.Master.QueryRowContext(
		ctx, query, rUrl.ShortURL, rUrl.UserAgent, rUrl.Device, rUrl.OS, rUrl.Browser, rUrl.IP,
	).Scan(&rUrl.ID)
	if err != nil {
		return "", fmt.Errorf("insert analytics row into db: %w", err)
	}

	return rUrl.ID, nil
}
