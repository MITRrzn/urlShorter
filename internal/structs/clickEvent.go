package structs

import "time"

type ClickEvent struct {
	EventID   string    `json:"event_id"`
	LinkID    int64     `json:"link_id"`
	ShortCode string    `json:"short_code"`
	ClickedAt time.Time `json:"clicked_at"`
	Referer   string    `json:"referer"`
	UserAgent string    `json:"user_agent"`
}
