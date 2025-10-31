package routes

import (
	"net/http"

	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/brij-812/url-shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, urlHandler *handlers.URLHandler, userHandler *handlers.UserHandler) {
	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Public routes
	r.Post("/signup", userHandler.Signup)
	r.Post("/login", userHandler.Login)

	// Protected routes (require valid JWT)
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuth)
		protected.Use(middleware.RateLimit) // 🔹 apply Redis rate limiter here

		protected.Post("/shorten", urlHandler.ShortenURL)
		protected.Get("/metrics", urlHandler.GetMetrics)
		protected.Get("/all", urlHandler.GetAllUserURLs)
	})

	// Public redirect (no auth or rate limit yet)
	r.Get("/{shortCode}", urlHandler.RedirectURL)
}
