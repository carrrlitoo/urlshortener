// handlers/handler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"urlshortener/database"
	"urlshortener/models"
	"urlshortener/service"
	"urlshortener/validation"

	"github.com/go-chi/chi/v5"
)

func HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}
}

func Shortener(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var longUrl models.RequestURL
		var shortCode string

		err := json.NewDecoder(r.Body).Decode(&longUrl)
		if err != nil {
			http.Error(w, "Invalid request json", http.StatusBadRequest)
			return
		}

		if longUrl.URL == "" || !validation.IsValidURL(longUrl.URL) {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		shortCode, err = database.GetShortCodeByURL(db, longUrl.URL)
		if err == nil && shortCode != "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.ResponseURL{ShortCode: shortCode})
			return
		}

		var shortCodeExists bool
		for {
			shortCode = service.GenerateShortCode()
			shortCodeExists, err = database.ShortCodeExists(db, shortCode)
			if err != nil {
				http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
				return
			}
			if !shortCodeExists {
				break
			}
		}

		fmt.Println("Generated short code:", shortCode)

		err = database.AddURLtoDB(db, longUrl.URL, shortCode)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(models.ResponseURL{ShortCode: shortCode})
	}
}

func GetAllURLs(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := database.GetAllURLsFromDB(db)
		if err != nil {
			http.Error(w, "Failed to retrieve URLs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urls)
	}
}

func RedirectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "shortCode")
		originalURL, err := database.GetURLbyShortCode(db, shortCode)
		if err != nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		err = database.IncrementClickCount(db, shortCode)
		if err != nil {
			http.Error(w, "Failed to update click count", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

func GetStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "shortCode")

		clicks, err := database.GetClickCount(db, shortCode)
		if err != nil {
			http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clicks)
	}
}

func DeleteURLbyShortCode(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestDeleteURL

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request json", http.StatusBadRequest)
			return
		}

		exists, err := database.ShortCodeExists(db, req.ShortCode)
		if err != nil {
			http.Error(w, "Failed to check short code existence", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Short code does not exist", http.StatusNotFound)
			return
		}

		err = database.DeleteURLByShortCode(db, req.ShortCode)
		if err != nil {
			http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "URL deleted successfully"})
	}
}
