# HyperLinkOS – Production-Grade URL Shortener (Go + PostgreSQL + Redis)

HyperLinkOS is a production-oriented URL shortening service built in Go, designed with clean architecture principles, robust authentication, caching, analytics, and full extensibility for modern deployments.  
The project supports PostgreSQL and in-memory modes, making it suitable for local development, testing, and scalable production deployments.

---

## Table of Contents

1. Overview  
2. Architecture  
3. Features  
4. Tech Stack  
5. Folder Structure  
6. Configuration  
7. Database Schema  
8. Redis Usage  
9. API Endpoints  
10. Core Workflows  
11. Running Locally  
12. Docker Usage  
13. Testing  
14. Roadmap  
15. License

---

## 1. Overview

HyperLinkOS provides a full backend service for creating, storing, and managing short URLs with per-user isolation, JWT authentication, domain analytics, and Redis caching.  
The backend exposes REST APIs, integrates cleanly with a frontend, and is prepared for containerized hosting environments.

---

## 2. Architecture

The system uses layered architecture:

- cmd/server: Application entry point, router setup, dependency initialization  
- internal/handlers: HTTP handlers  
- internal/repository: Database abstractions  
- internal/middleware: Authentication, rate limiting, logging  
- internal/utils: Helpers such as JWT, hashing, validators  
- configs: Koanf-based configuration files  

Key characteristics:

- Chi router with middleware chaining  
- JWT-based request authentication  
- Clean separation between business logic and persistence  
- Redis caching layer for hot paths (redirect and analytics)  
- Configuration handled by Koanf from YAML, env vars, or flags  

---

## 3. Features

Completed:

- Core backend with Koanf configuration  
- Signup and login with bcrypt + JWT  
- JWT middleware injecting user_id into context  
- URL shortening with per-user scoping  
- Public redirect endpoint  
- Domain-frequency metrics  
- Fetch all URLs owned by the user  
- Memory repository for test mode  
- Comprehensive unit and handler tests  
- Postgres-backed repository

Pending:

- Redis integration for caching and rate limiting  
- Per-user or per-IP rate limiting  
- URL expiration and TTL  
- URL deletion and editing  
- Click analytics (hits, referrer, timestamps, IP tracking)  
- Docker Compose environment with Redis and Postgres  
- Next.js frontend integration  
- CI/CD automation  
- Prometheus-style observability  

---

## 4. Tech Stack

Backend:

- Go 1.22+  
- Chi Router  
- PostgreSQL  
- Redis  
- Koanf configuration  
- bcrypt hashing  
- JSON Web Tokens  
- Docker (planned)  

Frontend (optional):

- Next.js 14  
- TypeScript  
- Typed client generated from OpenAPI  

---

## 5. Folder Structure

```
.
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── handlers
│   ├── repository
│   ├── middleware
│   ├── utils
│   └── models
├── configs
│   └── config.yaml
├── docker
│   ├── Dockerfile
│   └── docker-compose.yaml (planned)
└── README.md
```

---

## 6. Configuration

Configuration is loaded via Koanf using the following sources:

- YAML configuration files  
- Environment variables  
- CLI flags  

Key configuration values include:

- PostgreSQL DSN  
- Redis address  
- JWT secret and expiration interval  
- Application port  
- Mode: postgres or memory  

---

## 7. Database Schema

### links table

```
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ DEFAULT NULL
);
```

### domain_counts table

```
CREATE TABLE domain_counts (
    user_id INT NOT NULL,
    domain TEXT NOT NULL,
    count INT NOT NULL DEFAULT 1,
    PRIMARY KEY (user_id, domain)
);
```

Note: expires_at was added manually during development and should be included in future migrations.

---

## 8. Redis Usage

Redis caching improves performance for high-traffic systems.

Planned key usage:

- redirect:{code} — cache of long URLs  
- metrics:topdomains:{user_id} — cached domain analytics  
- rate:{user_id} — sliding-window rate limiting  

TTL values are configurable.

---

## 9. API Endpoints

### Public Endpoints

| Method | Path | Description |
|--------|-----------|-------------------------------|
| GET | /health | Health status |
| POST | /signup | User registration |
| POST | /login | User authentication |
| GET | /{code} | Redirect short code |

### Protected Endpoints (JWT Required)

| Method | Path | Description |
|--------|-----------|-----------------------------|
| POST | /shorten | Create new short URL |
| GET | /metrics | Domain-frequency metrics |
| GET | /all | Fetch all URLs of the user |
| DELETE | /url/{code} | Delete specific short URL |

Middleware applied:

1. JWTAuth  
2. RateLimit (planned)  

---

## 10. Core Workflows

### Shorten URL

1. User authenticates via JWT  
2. Handler extracts user_id  
3. URL is validated and normalized  
4. Short code is generated  
5. Entry stored in Postgres  
6. Domain metrics updated  
7. Metrics cache invalidated  

### Redirect

1. Check Redis cache  
2. Fallback to Postgres  
3. Verify expiry time  
4. Cache long URL if valid  
5. Respond with HTTP 301  

### Metrics

1. Attempt Redis lookup  
2. Query domain_counts if cache miss  
3. Cache results  
4. Return domain-to-count map  

---

## 11. Running Locally (Without Docker)

Prerequisites:

- PostgreSQL installed  
- Go 1.22+  

Setup:

```
go mod tidy
export APP_ENV=dev
go run cmd/server/main.go
```

Set up schema manually:

```
psql -d hyperlinkos -f schema.sql
```

---

## 12. Docker Usage

A complete Docker Compose environment will orchestrate:

- API service  
- PostgreSQL  
- Redis  
- Optional Next.js frontend  

Once completed:

```
docker-compose up --build
```

---

## 13. Testing

Tests include:

- Handler tests  
- Repository tests in memory mode  
- Metrics logic tests  

Run all:

```
go test ./...
```

---

## 14. Roadmap

- Redis-backed cache for redirects and metrics  
- Sliding-window rate limiting  
- TTL expiry and cleanup worker  
- URL update and deletion  
- Per-link click tracking  
- Next.js frontend integration  
- OpenAPI specification and typed client  
- CI/CD automation  
- Observability via Prometheus  

---

## 15. License

MIT License.
