// models/url.go

package models

type URL struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	Clicks      int    `json:"clicks"`
	CreatedAt   string `json:"created_at"`
}

type RequestURL struct {
	URL string `json:"url"`
}

type ResponseURL struct {
	ShortCode string `json:"short_code"`
}

type ResponseStats struct {
	Clicks int `json:"clicks"`
}

type RequestDeleteURL struct {
	ShortCode string `json:"short_code"`
}
