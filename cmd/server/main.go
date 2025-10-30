package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/brij-812/url-shortener/internal/cache"
	"github.com/brij-812/url-shortener/internal/config"
	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/brij-812/url-shortener/internal/middleware"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/routes"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	// 🔹 Load configuration
	cfg := config.LoadConfig()

	// 🔹 Initialize JWT secret for middleware (critical!)
	middleware.InitJWTSecret(cfg.JWT.Secret)

	// 🔹 Connect to Postgres
	db, err := sql.Open("postgres", "host="+cfg.Database.Host+
		" port="+cfg.Database.Port+
		" user="+cfg.Database.User+
		" password="+cfg.Database.Password+
		" dbname="+cfg.Database.Name+
		" sslmode="+cfg.Database.SSLMode)
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}
	defer db.Close()

	// 🔹 Initialize Redis (combine host + port)
	redisAddr := cfg.Redis.Host + ":" + cfg.Redis.Port
	cache.InitRedis(redisAddr, cfg.Redis.Password, cfg.Redis.DB)

	// 🔹 Initialize repository and handlers
	repo := repository.NewPostgresRepo(db)
	urlHandler := handlers.NewURLHandler(repo)
	userHandler := handlers.NewUserHandler(
		db,
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTokenExpiryMinutes,
	)

	// 🔹 Setup router
	r := chi.NewRouter()
	routes.RegisterRoutes(r, urlHandler, userHandler)

	// 🔹 Start server
	log.Printf("🚀 Server running on :%s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
