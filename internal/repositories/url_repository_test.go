package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"mini-url-shortener/internal/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	return db, mock
}

func TestCreateShortCode_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	url := &models.URL{OriginalURL: "https://test.test", ShortCode: "abc123"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO shorten_urls (original_url, short_code) VALUES (?, ?)")).
		WithArgs(url.OriginalURL, url.ShortCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateShortCode(context.Background(), url)
	assert.NoError(t, err)
}

func TestCreateShortCode_Duplicate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	url := &models.URL{OriginalURL: "https://test.test", ShortCode: "dupl1c4t3"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO shorten_urls (original_url, short_code) VALUES (?, ?)")).
		WithArgs(url.OriginalURL, url.ShortCode).
		WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))

	err := repo.CreateShortCode(context.Background(), url)
	assert.Error(t, err)
}

func TestFindByShortCode_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	shortCode := "abc123"
	expectedURL := models.URL{ID: 1, OriginalURL: "https://test.test"}

	mock.ExpectQuery("SELECT id, original_url FROM shorten_urls WHERE short_code = ?").
		WithArgs(shortCode).
		WillReturnRows(sqlmock.NewRows([]string{"id", "original_url"}).
			AddRow(expectedURL.ID, expectedURL.OriginalURL))

	url, err := repo.FindByShortCode(context.Background(), shortCode)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, expectedURL.ID, url.ID)
	assert.Equal(t, expectedURL.OriginalURL, url.OriginalURL)
}

func TestFindByShortCode_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	shortCode := "nf3211"
	mock.ExpectQuery("SELECT id, original_url FROM shorten_urls WHERE short_code = ?").
		WithArgs(shortCode).
		WillReturnError(sql.ErrNoRows)

	url, err := repo.FindByShortCode(context.Background(), shortCode)
	assert.Error(t, err)
	assert.Nil(t, url)
	assert.EqualError(t, err, fmt.Sprintf("short code '%s' not found", shortCode))
}

func TestUpdateClickByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE shorten_urls SET clicks = clicks + 1 WHERE id = ?")).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateClickByID(context.Background(), 999)
	assert.NoError(t, err)
}

func TestUpdateClickByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE shorten_urls SET clicks = clicks + 1 WHERE id = ?")).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateClickByID(context.Background(), 999)
	assert.Error(t, err)
	assert.EqualError(t, err, "no row updated")
}
