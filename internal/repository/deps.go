package repository

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

// querier методы операций в БД
type querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
