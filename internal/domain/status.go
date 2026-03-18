package domain

import "time"

type Status struct {
	SiteID     string    `json:"site_id"`
	IsUp       bool      `json:"is_up"`
	StatusCode int       `json:"status_code,omitempty"`
	Error      string    `json:"error,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
}