package repositories

import (
	"context"
	"database/sql"
	"mini-url-shortener/internal/models"
)

type URLRepository interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
	CreateShortCode(ctx context.Context, url *models.URL) error
}

type dbExecutor interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
