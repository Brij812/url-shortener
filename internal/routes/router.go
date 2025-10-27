package routes

import (
	"net/http"

	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes handles URL-related routes
func RegisterRoutes(r chi.Router, h *handlers.URLHandler) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Post("/shorten", h.ShortenURL)
	r.Get("/{code}", h.RedirectURL)
	r.Get("/metrics", h.GetMetrics)
}

// RegisterUserRoutes handles authentication routes
func RegisterUserRoutes(r chi.Router, uh *handlers.UserHandler) {
	r.Post("/signup", uh.Signup)
	r.Post("/login", uh.Login)
}
