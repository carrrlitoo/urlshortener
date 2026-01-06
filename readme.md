# URL Shortener

A URL shortening service built with Go and PostgreSQL. Similar to bit.ly, it converts long URLs into short, easy-to-share links and tracks click statistics.

## Features

- Shorten long URLs to compact codes
- Automatic redirect from short codes to original URLs
- Click tracking and statistics
- Delete shortened URLs
- Duplicate URL detection (returns existing code)
- Input validation
- Collision-resistant short code generation

## Technologies

- Go 1.25.4
- PostgreSQL 16
- Chi Router
- Docker & Docker Compose

## API Endpoints
```
POST   /shorten              - Create short URL
DELETE /shorten              - Delete URL by short code
GET    /{shortCode}          - Redirect to original URL
GET    /stats/{shortCode}    - Get click statistics
GET    /shorten/all          - Get all URLs (debug endpoint)
```

## Project Structure
```
urlshortener/
├── config/          - Configuration (.env handling)
├── database/        - PostgreSQL operations
├── handlers/        - HTTP handlers
├── models/          - Data structures
├── service/         - Business logic (code generation)
├── validation/      - Input validation
├── main.go          - Entry point
├── .env             - Environment variables
├── Dockerfile       - Application Docker image
├── docker-compose.yml - Container orchestration
└── init.sql         - Database schema
```

## Database Schema
```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    clicks INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Running the Project

### Using Docker (recommended)
```bash
docker-compose up --build
```

Application will be available at `http://localhost:8080`

### Local Development

1. Install PostgreSQL and create database
2. Create `.env` file with connection settings
3. Run application:
```bash
go run main.go
```

## Request Examples

### Shorten URL
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com/very/long/url"}'
```

Response:
```json
{"short_code": "abc123"}
```

### Access shortened URL
```bash
curl -L http://localhost:8080/abc123
```

Redirects to original URL.

### Get statistics
```bash
curl http://localhost:8080/stats/abc123
```

Response:
```json
{"clicks": 5}
```

### Delete shortened URL
```bash
curl -X DELETE http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"short_code": "abc123"}'
```

Response:
```json
{"message": "URL deleted successfully"}
```

### Local Development

1. Install PostgreSQL and create database
2. Copy `.env.example` to `.env` and fill in your database credentials:
```bash
   cp .env.example .env
```
3. Run application:
```bash
   go run main.go
```

## Features Implementation

- **URL Validation**: Validates scheme (http/https) and host presence
- **Duplicate Detection**: Returns existing short code for duplicate URLs
- **Collision Handling**: Regenerates code if collision detected
- **Click Tracking**: Increments counter on each redirect
- **Deletion**: Remove shortened URLs with validation
- **Error Handling**: Proper HTTP status codes for all scenarios