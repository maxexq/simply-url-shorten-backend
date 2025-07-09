package shortenrepo

import (
	"database/sql"

	model "url-shortener/internal/model/shorten-model"
)

type URLRepository interface {
	SaveURL(url *model.URL) error
	GetOriginalURL(shortCode string) (string, error)
}

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &urlRepository{db}
}

func (r *urlRepository) SaveURL(url *model.URL) error {
	_, err := r.db.Exec(
		"INSERT INTO urls (id, original_url, short_code) VALUES (?, ?, ?)",
		url.ID, url.OriginalURL, url.ShortCode,
	)
	return err
}

func (r *urlRepository) GetOriginalURL(shortCode string) (string, error) {
	var originalURL string
	err := r.db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", shortCode).Scan(&originalURL)
	return originalURL, err
}
