// urlshortener/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"urlshortener/config"
	"urlshortener/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	r := chi.NewRouter()

	connStr := config.GetDBConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка при проверке подключения к базе данных:", err)
		return
	}
	fmt.Println("Успешное подключение к базе данных! Successfully connected to the database!")

	r.Use(middleware.Logger)

	r.Get("/hello", handlers.HelloHandler())
	r.Post("/shorten", handlers.Shortener(db))
	r.Delete("/shorten", handlers.DeleteURLbyShortCode(db))
	r.Get("/shorten/all", handlers.GetAllURLs(db))

	r.Get("/{shortCode}", handlers.RedirectHandler(db))

	r.Get("/stats/{shortCode}", handlers.GetStatsHandler(db))

	fmt.Println("Server is running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
