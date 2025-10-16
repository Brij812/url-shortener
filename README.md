# URL Shortener

A simple URL shortener service written in Go. It allows you to shorten URLs, redirect short URLs to their original URLs, and view top accessed domains (metrics).

## Features

* Shorten URLs
* Redirect short URLs
* Track top domains accessed
* Thread-safe in-memory storage
* Docker-ready

## Project Structure

url-shortener/
├── cmd/
│   └── server/
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── utils/
│   └── routes/
├── Dockerfile
├── go.mod
└── README.md

## Local Setup

Clone the repo:

```bash
git clone https://github.com/Brij812/url-shortener
cd url-shortener
```

Install dependencies:

```bash
go mod tidy
```

Run the server:

```bash
go run ./cmd/server
```

Server runs on: `http://localhost:8080`

## API Endpoints

Shorten URL:

```bash
curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"url":"https://www.example.com"}'
```

Redirect short URL:

```bash
curl http://localhost:8080/<code>
```

Metrics (top domains):

```bash
curl http://localhost:8080/metrics
```

## Testing

Run all tests:

```bash
go test ./... -v
```

## Docker Usage

Build image:

```bash
docker build -t url-shortener .
```

Run container:

```bash
docker run -p 8080:8080 url-shortener
```

## Git Ignore

```
/bin/
/obj/
/vendor/
venv/
*.vscode/
*.idea/
.DS_Store
Thumbs.db
```
