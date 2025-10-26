package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/brij-812/url-shortener/internal/config"
	"github.com/brij-812/url-shortener/internal/database"
	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewPostgresDB(cfg)
	defer db.Close()

	migrateFlag := flag.String("migrate", "", "Run DB migrations: up, down, or version")
	flag.Parse()

	if *migrateFlag != "" {
		database.RunMigrations(db, cfg, *migrateFlag)
		return
	}

	repo := repository.NewPostgresRepo(db)
	handler := handlers.NewURLHandler(repo)

	r := chi.NewRouter()
	routes.RegisterRoutes(r, handler)

	log.Printf("🚀 Server running on port %s\n", cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, r)
}
