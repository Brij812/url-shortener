package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/brij-812/url-shortener/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies or rolls back database migrations based on the given direction ("up", "down", "version").
func RunMigrations(db *sql.DB, cfg *config.Config, direction string) {
	// Initialize migration driver using existing DB connection
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Migration driver init failed: %v", err)
	}

	// Dynamically resolve absolute path to migrations folder
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Failed to get working directory: %v", err)
	}

	absPath := filepath.Join(wd, "migrations")

	// On Windows, use file://C:/path instead of file:///C:/path
	migrationPath := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))
	log.Printf("📂 Using migration path: %s", migrationPath)

	// Load migrations
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		cfg.Database.Name,
		driver,
	)
	if err != nil {
		log.Fatalf("❌ Failed to load migrations: %v", err)
	}

	// Handle migration commands
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("❌ Migration up failed: %v", err)
		}
		log.Println("✅ Migrations applied successfully")

	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("❌ Migration down failed: %v", err)
		}
		log.Println("⬅️  Rolled back one migration")

	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				log.Println("ℹ️  No migrations applied yet")
				return
			}
			log.Fatalf("❌ Failed to get migration version: %v", err)
		}
		log.Printf("📦 Current DB version: %d (dirty: %v)\n", v, dirty)

	default:
		log.Fatalf("❌ Invalid migration command: %s (use 'up', 'down', or 'version')", direction)
	}
}
