package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/brij-812/url-shortener/internal/config"
	"github.com/brij-812/url-shortener/internal/database"
	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/brij-812/url-shortener/internal/middleware"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load config.yaml + env
	cfg := config.LoadConfig()

	// Connect to Postgres
	db := database.NewPostgresDB(cfg)
	defer db.Close()
	log.Println("✅ Connected to Postgres successfully")

	// Handle migrations
	migrateFlag := flag.String("migrate", "", "Run DB migrations: up, down, or version")
	flag.Parse()
	if *migrateFlag != "" {
		database.RunMigrations(db, cfg, *migrateFlag)
		return
	}

	// Initialize repo and handlers
	repo := repository.NewPostgresRepo(db)
	urlHandler := handlers.NewURLHandler(repo)
	userHandler := handlers.NewUserHandler(
		db,
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTokenExpiryMinutes,
	)

	// Initialize router
	r := chi.NewRouter()

	// Inject config secret into middleware (global secret)
	middleware.InitJWTSecret(cfg.JWT.Secret)

	// Register routes
	routes.RegisterRoutes(r, urlHandler, userHandler)

	// Start server
	log.Printf("🚀 Server running on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
