package models

import "time"

// Request body for shortening a URL
type ShortenRequest struct {
	URL        string `json:"url"`
	ExpiryDays int    `json:"expiry_days,omitempty"`
}

// Response body for a shortened URL
type ShortenResponse struct {
	ShortURL  string     `json:"short_url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
