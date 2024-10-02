package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Storage реализует методы работы с базой
type Storage struct {
	db *sqlx.DB
}

// New создает хранилище
func New(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Storage) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}

func (s *Storage) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

// Check проверяет доступность БД
func (s *Storage) Check() (interface{}, error) {
	return s.db.Stats(), s.db.Ping()
}
