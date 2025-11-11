package routes

import (
	"net/http"

	"github.com/brij-812/HyperLinkOS/internal/handlers"
	"github.com/brij-812/HyperLinkOS/internal/middleware"
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
	r.Post("/logout", userHandler.Logout)

	// 🔹 Protected APIs (require JWT)
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuth)

		// 🧠 Apply rate limiting *only* on /shorten
		protected.Group(func(limited chi.Router) {
			limited.Use(middleware.RateLimit)
			limited.Post("/shorten", urlHandler.ShortenURL)
		})

		// Normal protected endpoints (no rate limit)
		protected.Get("/metrics", urlHandler.GetMetrics)
		protected.Get("/all", urlHandler.GetAllUserURLs)
		protected.Delete("/url/{code}", urlHandler.DeleteURL)
	})

	// 🔹 Public redirect route
	r.Get("/{shortCode:[A-Za-z0-9_-]{4,12}}", urlHandler.RedirectURL)
}
