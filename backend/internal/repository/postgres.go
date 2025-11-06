package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/brij-812/HyperLinkOS/internal/cache"
)

// PostgresRepo stores data in Postgres instead of memory.
type PostgresRepo struct {
	db *sql.DB
}

// Constructor
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// 🧩 Extracts clean domain (removes www.)
func extractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return ""
	}
	host := strings.ToLower(u.Host)
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}
	return host
}

// Save inserts a new URL–code pair associated with a user.
// Supports optional expiry (TTL). If expiresAt is nil, link never expires.
func (r *PostgresRepo) Save(u, code string, userID int, expiresAt *time.Time) {
	_, err := r.db.ExecContext(context.Background(), `
		INSERT INTO links (code, long_url, user_id, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (code) DO NOTHING
	`, code, u, userID, time.Now(), expiresAt)
	if err != nil {
		log.Printf("❌ Failed to save URL: %v", err)
		return
	}

	// Increment domain count (user-specific)
	domain := extractDomain(u)
	log.Printf("🧩 Extracted domain for %s = '%s'", u, domain)
	if domain != "" {
		log.Printf("🧠 DEBUG: Save() called for URL=%s userID=%d domain=%s", u, userID, domain)
		_, err = r.db.ExecContext(context.Background(), `
			INSERT INTO domain_counts (domain, user_id, count)
			VALUES ($1, $2, 1)
			ON CONFLICT (domain, user_id)
			DO UPDATE SET count = domain_counts.count + 1
		`, domain, userID)
		if err != nil {
			log.Printf("❌ INSERT domain_counts failed: %v", err)
		} else {
			log.Printf("✅ INSERT domain_counts succeeded for domain=%s userID=%d", domain, userID)
		}
	}

	// 🧹 Invalidate cached metrics for this user
	cache.Delete(fmt.Sprintf("metrics:topdomains:%d", userID))
}

// GetCode finds the short code for a given long URL
func (r *PostgresRepo) GetCode(u string) (string, bool) {
	var code string
	err := r.db.QueryRow(`SELECT code FROM links WHERE long_url=$1`, u).Scan(&code)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		log.Printf("❌ GetCode error: %v", err)
		return "", false
	}
	return code, true
}

// GetURL finds the original long URL for a given code (public).
// Automatically skips expired links.
func (r *PostgresRepo) GetURL(code string) (string, bool) {
	cacheKey := "shorturl:" + code

	// 1️⃣ check Redis cache first
	if cached, ok := cache.Get(cacheKey); ok {
		return cached, true
	}

	// 2️⃣ fallback to Postgres
	var u string
	var expiresAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT long_url, expires_at
		FROM links
		WHERE code = $1
	`, code).Scan(&u, &expiresAt)

	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		log.Printf("❌ GetURL error: %v", err)
		return "", false
	}

	// 🕓 Check expiry
	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		log.Printf("⚰️ Link %s expired at %v", code, expiresAt.Time)
		return "", false
	}

	// 3️⃣ cache result in Redis (set TTL to min(24h, remaining validity))
	ttl := 24 * time.Hour
	if expiresAt.Valid {
		remaining := time.Until(expiresAt.Time)
		if remaining > 0 && remaining < ttl {
			ttl = remaining
		}
	}
	cache.Set(cacheKey, u, ttl)

	return u, true
}

// GetTopDomains returns top N most frequently saved domains for a specific user
func (r *PostgresRepo) GetTopDomains(userID, n int) map[string]int {
	cacheKey := fmt.Sprintf("metrics:topdomains:%d", userID)

	// 1️⃣ Try Redis cache
	if cachedJSON, ok := cache.Get(cacheKey); ok {
		out := make(map[string]int)
		if err := json.Unmarshal([]byte(cachedJSON), &out); err == nil {
			return out
		}
	}

	// 2️⃣ Query DB if cache miss
	rows, err := r.db.Query(`
		SELECT domain, count FROM domain_counts
		WHERE user_id = $1
		ORDER BY count DESC
		LIMIT $2
	`, userID, n)
	if err != nil {
		log.Printf("❌ GetTopDomains error: %v", err)
		return map[string]int{}
	}
	defer rows.Close()

	out := make(map[string]int)
	for rows.Next() {
		var domain string
		var count int
		if err := rows.Scan(&domain, &count); err == nil {
			out[domain] = count
		}
	}

	// 3️⃣ Save to Redis for 10 minutes
	data, _ := json.Marshal(out)
	cache.Set(cacheKey, string(data), 10*time.Minute)

	return out
}

// IncrementDomainCount increases count for a given domain (user-specific)
func (r *PostgresRepo) IncrementDomainCount(u string, userID int) {
	domain := extractDomain(u)
	if domain == "" {
		return
	}
	_, err := r.db.Exec(`
		INSERT INTO domain_counts (domain, user_id, count)
		VALUES ($1, $2, 1)
		ON CONFLICT (domain, user_id)
		DO UPDATE SET count = domain_counts.count + 1
	`, domain, userID)
	if err != nil {
		log.Printf("❌ IncrementDomainCount error: %v", err)
	}

	cache.Delete(fmt.Sprintf("metrics:topdomains:%d", userID))
}

// GetAllURLsByUser returns all shortened URLs for a given user
func (r *PostgresRepo) GetAllURLsByUser(userID int) []map[string]string {
	rows, err := r.db.Query(`
		SELECT code, long_url, created_at, expires_at
		FROM links
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		log.Printf("❌ GetAllURLsByUser error: %v", err)
		return nil
	}
	defer rows.Close()

	var results []map[string]string
	for rows.Next() {
		var code, longURL string
		var createdAt time.Time
		var expiresAt sql.NullTime
		if err := rows.Scan(&code, &longURL, &createdAt, &expiresAt); err == nil {
			expiry := ""
			if expiresAt.Valid {
				expiry = expiresAt.Time.Format(time.RFC3339)
			}
			results = append(results, map[string]string{
				"short_url":  "http://localhost:8080/" + code,
				"long_url":   longURL,
				"created_at": createdAt.Format(time.RFC3339),
				"expires_at": expiry,
			})
		}
	}
	return results
}

// 🧹 Optional: CleanupExpiredLinks removes expired links from DB & cache.
func (r *PostgresRepo) CleanupExpiredLinks() {
	ctx := context.Background()
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM links
		WHERE expires_at IS NOT NULL AND expires_at < NOW()
	`)
	if err != nil {
		log.Printf("❌ CleanupExpiredLinks error: %v", err)
		return
	}
	rows, _ := res.RowsAffected()
	log.Printf("🧹 CleanupExpiredLinks removed %d expired links", rows)
}

func (r *PostgresRepo) DeleteLink(userID int, code string) bool {
	res, err := r.db.Exec(`
		DELETE FROM links
		WHERE code = $1 AND user_id = $2
	`, code, userID)
	if err != nil {
		log.Printf("❌ DeleteLink error: %v", err)
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false
	}

	cache.Delete("shorturl:" + code)
	cache.Delete(fmt.Sprintf("metrics:topdomains:%d", userID))

	return true
}
