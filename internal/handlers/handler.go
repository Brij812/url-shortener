package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brij-812/url-shortener/internal/models"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	Repo *repository.MemoryRepo
}

func NewURLHandler(repo *repository.MemoryRepo) *URLHandler {
	return &URLHandler{Repo: repo}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url field required", http.StatusBadRequest)
		return
	}

	if code, exists := h.Repo.GetCode(req.URL); exists {
		json.NewEncoder(w).Encode(models.ShortenResponse{
			ShortURL: "http://localhost:8080/" + code,
		})
		return
	}

	code := utils.GenerateShortCode(req.URL)
	h.Repo.Save(req.URL, code)

	resp := models.ShortenResponse{
		ShortURL: "http://localhost:8080/" + code,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "code missing", http.StatusBadRequest)
		return
	}

	url, exists := h.Repo.GetURL(code)
	if !exists {
		http.Error(w, "short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
