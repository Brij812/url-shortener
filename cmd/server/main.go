package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/brij-812/url-shortener/internal/cache"
	"github.com/brij-812/url-shortener/internal/config"
	"github.com/brij-812/url-shortener/internal/database"
	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/brij-812/url-shortener/internal/middleware"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/routes"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

// connectWithRetry tries to connect to Postgres several times before giving up
func connectWithRetry(dsn string, retries int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Connected to Postgres")
				return db, nil
			}
		}
		log.Printf("⏳ Waiting for DB... (%d/%d)", i+1, retries)
		time.Sleep(3 * time.Second)
	}
	return nil, err
}

func main() {
	// 🔹 Load configuration
	cfg := config.LoadConfig()

	// 🔹 Initialize JWT secret for middleware (critical!)
	middleware.InitJWTSecret(cfg.JWT.Secret)

	// 🔹 Build DSN
	dsn := "host=" + cfg.Database.Host +
		" port=" + cfg.Database.Port +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.Name +
		" sslmode=" + cfg.Database.SSLMode

	// 🔹 Connect to Postgres with retry logic
	db, err := connectWithRetry(dsn, 10)
	if err != nil {
		log.Fatalf("❌ DB connection failed after retries: %v", err)
	}
	defer db.Close()

	// ✅ Run migrations automatically
	database.RunMigrations(db, cfg, "up")

	// 🔹 Initialize Redis
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
