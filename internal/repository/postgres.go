package repository

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"time"
)

// PostgresRepo stores data in Postgres instead of memory.
type PostgresRepo struct {
	db *sql.DB
}

// Constructor
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Save inserts a new URL–code pair associated with a user.
func (r *PostgresRepo) Save(u, code string, userID int) {
	_, err := r.db.ExecContext(context.Background(), `
		INSERT INTO links (code, long_url, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (code) DO NOTHING
	`, code, u, userID, time.Now())
	if err != nil {
		log.Printf("❌ Failed to save URL: %v", err)
		return
	}

	// Increment domain count
	d, err := url.Parse(u)
	if err == nil && d.Host != "" {
		_, err = r.db.ExecContext(context.Background(), `
			INSERT INTO domain_counts (domain, count)
			VALUES ($1, 1)
			ON CONFLICT (domain)
			DO UPDATE SET count = domain_counts.count + 1
		`, d.Host)
		if err != nil {
			log.Printf("❌ Failed to update domain count: %v", err)
		}
	}
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

// GetURL finds the original long URL for a given code (public)
func (r *PostgresRepo) GetURL(code string) (string, bool) {
	var u string
	err := r.db.QueryRow(`SELECT long_url FROM links WHERE code=$1`, code).Scan(&u)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		log.Printf("❌ GetURL error: %v", err)
		return "", false
	}
	return u, true
}

// GetTopDomains returns top N most frequently saved domains
func (r *PostgresRepo) GetTopDomains(n int) map[string]int {
	rows, err := r.db.Query(`
		SELECT domain, count FROM domain_counts
		ORDER BY count DESC
		LIMIT $1
	`, n)
	if err != nil {
		log.Printf("❌ GetTopDomains error: %v", err)
		return nil
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
	return out
}

// IncrementDomainCount increases count for a given domain
func (r *PostgresRepo) IncrementDomainCount(u string) {
	d, err := url.Parse(u)
	if err != nil || d.Host == "" {
		return
	}
	_, err = r.db.Exec(`
		INSERT INTO domain_counts (domain, count)
		VALUES ($1, 1)
		ON CONFLICT (domain)
		DO UPDATE SET count = domain_counts.count + 1
	`, d.Host)
	if err != nil {
		log.Printf("❌ IncrementDomainCount error: %v", err)
	}
}

// GetAllURLsByUser returns all shortened URLs for a given user
func (r *PostgresRepo) GetAllURLsByUser(userID int) []map[string]string {
	rows, err := r.db.Query(`
		SELECT code, long_url, created_at
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
		if err := rows.Scan(&code, &longURL, &createdAt); err == nil {
			results = append(results, map[string]string{
				"short_url":  "http://localhost:8080/" + code,
				"long_url":   longURL,
				"created_at": createdAt.Format(time.RFC3339),
			})
		}
	}
	return results
}
