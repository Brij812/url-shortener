package main

import (
	"log"
	"net/http"

	"github.com/brij-812/url-shortener/internal/handler"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo := repository.NewMemoryRepo()
	handler := handler.NewURLHandler(repo)

	r := chi.NewRouter()
	routes.RegisterRoutes(r, handler)

	log.Println("🚀 Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
