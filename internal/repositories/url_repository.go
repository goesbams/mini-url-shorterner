package repositories

import (
	"database/sql"
)

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &urlRepository{
		db: db,
	}
}
