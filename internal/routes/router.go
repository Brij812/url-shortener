package routes

import (
	"net/http"

	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/brij-812/url-shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes wires up all API endpoints.
func RegisterRoutes(r chi.Router, urlHandler *handlers.URLHandler, userHandler *handlers.UserHandler) {
	// 🔹 Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 🔹 Public authentication routes
	r.Post("/signup", userHandler.Signup)
	r.Post("/login", userHandler.Login)

	// 🔹 Protected APIs (require JWT + RateLimit)
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuth)
		protected.Use(middleware.RateLimit)

		protected.Post("/shorten", urlHandler.ShortenURL)
		protected.Get("/metrics", urlHandler.GetMetrics)
		protected.Get("/all", urlHandler.GetAllUserURLs)
		protected.Delete("/url/{code}", urlHandler.DeleteURL)
	})

	r.Get("/{shortCode:[A-Za-z0-9_-]{4,12}}", urlHandler.RedirectURL)
}
