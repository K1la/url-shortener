package model

import "time"

type URL struct {
	ID        string    `json:"id"`
	URL       string    `json:"url,omitempty"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}

type RedirectClicks struct {
	ID        string    `json:"id"`
	ShortURL  string    `json:"short_url"`
	UserAgent string    `json:"user_agent"`
	Device    string    `json:"device"`
	OS        string    `json:"os"`
	Browser   string    `json:"browser"`
	CreatedAt time.Time `json:"created_at"`
}
