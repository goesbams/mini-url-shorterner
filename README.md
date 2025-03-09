# Mini URL Shortener API

This is a mini URL shortener API built with Golang, featuring unit tests and mock services.

## Entity Relationship Diagram (ERD)

```plaintext
+------------------+
|     URLs        |
+------------------+
| id              |
| original_url    |
| short_code      |
| created_at      |
+------------------+
```

## Sequence Diagram

```plaintext
User ->> Server: POST /shorten {original_url}
Server ->> DB: Store short URL
Server ->> User: Returns shortened URL

User ->> Server: GET /{short_code}
Server ->> DB: Retrieve original URL
Server ->> DB: Update clicked counts
Server ->> User: Redirect to original URL
```

## Project Structure

```plaintext
/mini-url-shortener

│
├───cmd
│   └───api
│           main.go
│
├───config
│       config.go
│       config.local.yml
│
├───internal
│   ├───database
│   │       database.go
│   │       database_test.go
│   │
│   ├───handlers
│   │       handlers.go
│   │       handlers_test.go
│   │
│   ├───helpers
│   │       helpers.go
│   │       helpers_test.go
│   │
│   ├───models
│   │       models.go
│   │
│   ├───repositories
│   │       repositories.go
│   │       url_repository.go
│   │       url_repository_test.go
│   │
│   ├───routes
│   │       routes.go
│   │
│   ├───server
│   │       server.go
│   │
│   └───services
│           services.go
│           url_service.go
│           url_service_test.go
│
├───migrations
│       1_create_shorten_urls_table.down.sql
│       1_create_shorten_urls_table.up.sql
│   go.mod
│   go.sum
│   Makefile
│   README.MD
```

## Installation & Setup

1. Clone this repository:
   ```sh
   git clone https://github.com/yourname/mini-url-shortener.git
   cd mini-url-shortener
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Setup database:
   ```sh
   make setup-db
   ```
4. Run migrations:
   ```sh
   make migrate-up
   ```
5. Run the project:
   ```sh
   make run-api
   ```

### Other Commands:
- Run tests: `make test`
- Clean up migrations: `make clean`

## API Endpoints

### Shorten URL
**Request:**
```http
POST /shorten
Content-Type: application/json
{
  "original_url": "https://example.com"
}
```
**Response:**
```json
{
  "short_code": "abc123"
}
```

### Redirect to Original URL
**Request:**
```http
GET /abc123
```
**Response:**
Redirects to `https://example.com`

### Health Check
**Request:**
```http
GET /ping
```
**Response:**
```text
pong
```

## Testing
Run unit tests with:
```sh
make test
```

