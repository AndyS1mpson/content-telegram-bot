package repository

import "time"

// Pin описывает контент
type Pin struct {
	ID        int64     `db:"id"`
	URL       string    `db:"url"`
	Type      string    `db:"type"`
	Status    int64     `db:"status"`
	Channel   string    `db:"channel"`
	CreatedAt time.Time `db:"created_at"`
}
