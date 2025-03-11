package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"mini-url-shortener/internal/models"
)

type urlRepository struct {
	db     *sql.DB
	tx     *sql.Tx
	dbExec dbExecutor
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &urlRepository{
		db:     db,
		dbExec: db,
	}
}

func (r *urlRepository) Begin(ctx context.Context) error {
	if r.tx != nil {
		return fmt.Errorf("transaction already in progress")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	r.tx = tx
	r.dbExec = tx

	return nil
}

func (r *urlRepository) Commit() error {
	if r.tx == nil {
		return fmt.Errorf("no active transaction to commit")
	}

	if err := r.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.tx = nil
	r.dbExec = r.db

	return nil
}

func (r *urlRepository) Rollback() error {
	if r.tx == nil {
		return fmt.Errorf("no active transaction to rollback")
	}

	if err := r.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	r.tx = nil
	r.dbExec = r.db

	return nil
}

func (r *urlRepository) CreateShortCode(ctx context.Context, url *models.URL) error {
	query := "INSERT INTO urls (original_url, short_code) VALUES (:original_url, :short_code)"
	_, err := r.dbExec.ExecContext(ctx, query, sql.Named("original_url", url.OriginalURL), sql.Named("short_code", url.ShortCode))
	if err != nil {
		return fmt.Errorf("failed to create short code: %w", err)
	}

	return nil
}
