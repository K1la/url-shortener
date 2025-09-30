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
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

type SummaryOfAnalytics struct {
	ShortUrl    string         `json:"short_url"`
	TotalClicks int            `json:"total_clicks"`
	Daily       map[string]int `json:"daily"`      // clicks per day
	Monthly     map[string]int `json:"monthly"`    // clicks per month
	UserAgent   map[string]int `json:"user_agent"` // clicks per User_Agent
}
