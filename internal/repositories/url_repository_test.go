package repositories

import (
	"context"
	"fmt"
	"mini-url-shortener/internal/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortCode_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	url := &models.URL{OriginalURL: "https://test.test", ShortCode: "abc123"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO shorten_urls (original_url, short_code) VALUES (:original_url, :short_code)")).
		WithArgs(url.OriginalURL, url.ShortCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateShortCode(context.Background(), url)
	assert.NoError(t, err)
}

func TestCreateShortCode_Duplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &urlRepository{
		db:     db,
		dbExec: db,
	}

	url := &models.URL{OriginalURL: "https://test.test", ShortCode: "dupl1c4t3"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO shorten_urls (original_url, short_code) VALUES (:original_url, :short_code)")).
		WithArgs(url.OriginalURL, url.ShortCode).
		WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))

	err = repo.CreateShortCode(context.Background(), url)
	assert.Error(t, err)
}
