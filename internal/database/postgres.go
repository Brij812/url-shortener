package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/brij-812/url-shortener/internal/config"
	_ "github.com/lib/pq" // Postgres driver
)

func NewPostgresDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping Postgres: %v", err)
	}

	log.Println("✅ Connected to Postgres successfully")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)

	return db
}
