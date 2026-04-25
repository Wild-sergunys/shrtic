package model

import "time"

type Link struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"-"` // указатель чтобы различать юзера
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}
