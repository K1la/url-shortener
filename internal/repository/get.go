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

func (r *Repository) GetShortURL(ctx context.Context, shortURL string) (*model.URL, error) {
	query := `
		SELECT id, url, short_url
		FROM urls
		WHERE short_url = $1
	`

	var url model.URL
	err := r.db.Master.QueryRowContext(ctx, query, shortURL).
		Scan(&url.ID, &url.URL, &url.ShortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrShortURLNotFound
		}

		return nil, fmt.Errorf("get short url: %w", err)
	}
	return &url, nil
}

func (r *Repository) CountClicks(ctx context.Context, shortURL string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM analytics 
		WHERE short_url = $1;
	`
	var count int
	err := r.db.Master.QueryRowContext(ctx, query, shortURL).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count clicks: %w", err)
	}
	return count, nil
}

func (r *Repository) CountClicksByDay(ctx context.Context, shortURL string) (map[string]int, error) {
	query := `
		SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS day, COUNT(*)
		FROM analytics
		WHERE short_url = $1
		GROUP BY day
		ORDER BY day DESC;
	`
	rows, err := r.db.Master.QueryContext(ctx, query, shortURL)
	if err != nil {
		return nil, fmt.Errorf("count query clicks by day: %w", err)
	}
	defer rows.Close()

	res := make(map[string]int)
	for rows.Next() {
		var day string
		var count int

		if err := rows.Scan(&day, &count); err != nil {
			return nil, fmt.Errorf("scan clicks by day: %w", err)
		}
		res[day] = count
	}

	return res, nil
}

func (r *Repository) CountClicksByMonth(ctx context.Context, shortURL string) (map[string]int, error) {
	query := `
		SELECT TO_CHAR(created_at, 'YYYY-MM') AS month, COUNT(*)
		FROM analytics
		WHERE short_url = $1
		GROUP BY month
		ORDER BY month DESC;
	`
	rows, err := r.db.Master.QueryContext(ctx, query, shortURL)
	if err != nil {
		return nil, fmt.Errorf("count query clicks by month: %w", err)
	}
	defer rows.Close()

	res := make(map[string]int)
	for rows.Next() {
		var month string
		var count int

		if err := rows.Scan(&month, &count); err != nil {
			return nil, fmt.Errorf("scan clicks by month: %w", err)
		}
		res[month] = count
	}

	return res, nil
}

func (r *Repository) CountClicksByUserAgent(ctx context.Context, shortURL string) (map[string]int, error) {
	query := `
		SELECT user_agent, COUNT(*) 
		FROM analytics
		WHERE short_url = $1
		GROUP BY user_agent
		ORDER BY COUNT(*) DESC;
	`

	rows, err := r.db.Master.QueryContext(ctx, query, shortURL)
	if err != nil {
		return nil, fmt.Errorf("query clicks by user agent: %w", err)
	}
	defer rows.Close()

	res := make(map[string]int)
	for rows.Next() {
		var ua string
		var count int

		if err := rows.Scan(&ua, &count); err != nil {
			return nil, fmt.Errorf("scan clicks by user agent: %w", err)
		}

		res[ua] = count
	}

	return res, nil
}
