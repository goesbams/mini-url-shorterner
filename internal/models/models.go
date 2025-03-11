package models

type URL struct {
	ID          int    `db:"id"`
	OriginalURL string `db:"original_url"`
	ShortCode   string `db:"short_code"`
	Clicks      int    `db:"clicks"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	DeletedAt   string `db:"deleted_at"`
}

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

type URLResponse struct {
	ShortCode string `json:"short_code"`
}
