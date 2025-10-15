package routes

import (
	"net/http"

	"github.com/brij-812/url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *handlers.URLHandler) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Post("/shorten", h.ShortenURL)
	r.Get("/{code}", h.RedirectURL) //
}
