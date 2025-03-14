package repositories

import (
	"context"
	"database/sql"
	"errors"
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
	query := "INSERT INTO shorten_urls (original_url, short_code) VALUES (?, ?)"
	_, err := r.dbExec.ExecContext(ctx, query, url.OriginalURL, url.ShortCode)
	if err != nil {
		return fmt.Errorf("failed to create short code: %w", err)
	}

	return nil
}

func (r *urlRepository) FindByShortCode(ctx context.Context, shortenCode string) (*models.URL, error) {
	var url models.URL

	err := r.dbExec.QueryRowContext(ctx, "SELECT id, original_url FROM shorten_urls WHERE short_code = ?", shortenCode).
		Scan(&url.ID, &url.OriginalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("short code '%s' not found", shortenCode)
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	return &url, nil
}

func (r *urlRepository) UpdateClickByID(ctx context.Context, id int) error {
	res, err := r.dbExec.ExecContext(ctx, "UPDATE shorten_urls SET clicks = clicks + 1 WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to update click count: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no row updated")
	}
	return nil
}
