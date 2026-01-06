// database/db.go

package database

import (
	"database/sql"
	"urlshortener/models"
)

func AddURLtoDB(db *sql.DB, originalURL, shortCode string) error {
	_, err := db.Exec("INSERT INTO urls (original_url, short_code, clicks, created_at) VALUES ($1, $2, 0, NOW())", originalURL, shortCode)
	return err
}

func GetAllURLsFromDB(db *sql.DB) ([]models.URL, error) {
	rows, err := db.Query("SELECT id, original_url, short_code, clicks, created_at FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		if err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.Clicks, &url.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func GetURLbyShortCode(db *sql.DB, shortCode string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalURL)
	return originalURL, err
}

func GetShortCodeByURL(db *sql.DB, URL string) (string, error) {
	var shortCode string
	err := db.QueryRow("SELECT short_code FROM urls WHERE original_url = $1", URL).Scan(&shortCode)
	return shortCode, err
}

func ShortCodeExists(db *sql.DB, shortCode string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)", shortCode).Scan(&exists)
	return exists, err
}

func IncrementClickCount(db *sql.DB, shortCode string) error {
	_, err := db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_code = $1", shortCode)
	return err
}

func GetClickCount(db *sql.DB, shortCode string) (models.ResponseStats, error) {
	var clicks models.ResponseStats
	err := db.QueryRow("SELECT clicks FROM urls WHERE short_code = $1", shortCode).Scan(&clicks.Clicks)

	return clicks, err
}

func DeleteURLByShortCode(db *sql.DB, shortCode string) error {
	_, err := db.Exec("DELETE FROM urls WHERE short_code = $1", shortCode)
	return err
}
