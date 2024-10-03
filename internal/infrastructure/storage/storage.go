package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	syncStorage = &Storage{}
	syncErr     error
	initOnce    = sync.Once{}
)

// Storage реализует методы работы с базой
type Storage struct {
	db *sqlx.DB
}

// New создает хранилище
func New(config Config) *Storage {
	initOnce.Do(func() {
		conn, err := sqlx.Open(
			"postgres",
			fmt.Sprintf(
				"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				config.Host,
				config.Port,
				config.User,
				config.Password,
				config.DB,
			))
		if err != nil {
			panic(fmt.Sprintf("init storage conn: %s", err))
		}

		conn.SetMaxOpenConns(10)
		conn.SetMaxIdleConns(5)
		conn.SetConnMaxLifetime(time.Minute * 15)
		syncStorage.db = conn
	})
	if syncErr != nil {
		panic(fmt.Errorf("failed to create db connection: %w", syncErr))
	}
	return syncStorage
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

func (s *Storage) Close() error {
	return s.db.Close()
}

// Check проверяет доступность БД
func (s *Storage) Check() (interface{}, error) {
	return s.db.Stats(), s.db.Ping()
}
