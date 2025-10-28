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

	// 🔹 Public redirect route (anyone can use short link)
	r.Get("/{shortCode}", urlHandler.RedirectURL)

	// Protected routes (require valid JWT)
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuth)
		protected.Post("/shorten", urlHandler.ShortenURL)
		protected.Get("/metrics", urlHandler.GetMetrics)
		protected.Get("/all", urlHandler.GetAllUserURLs)
	})
}
