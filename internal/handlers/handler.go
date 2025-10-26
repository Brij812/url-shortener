package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/brij-812/url-shortener/internal/models"
	"github.com/brij-812/url-shortener/internal/repository"
	"github.com/brij-812/url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	Repo repository.Repository // now uses the interface, not MemoryRepo
}

func NewURLHandler(repo repository.Repository) *URLHandler {
	return &URLHandler{Repo: repo}
}

func normalizeURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return raw
	}
	u.Fragment = ""
	if strings.HasSuffix(u.Path, "/") && u.Path != "/" {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}
	return u.String()
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

	req.URL = normalizeURL(req.URL)

	code, exists := h.Repo.GetCode(req.URL)
	if !exists {
		code = utils.GenerateShortCode(req.URL)
		h.Repo.Save(req.URL, code)
	} else {
		h.Repo.IncrementDomainCount(req.URL)
	}

	json.NewEncoder(w).Encode(models.ShortenResponse{
		ShortURL: "http://localhost:8080/" + code,
	})
}

func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "code missing", http.StatusBadRequest)
		return
	}

	longURL, exists := h.Repo.GetURL(code)
	if !exists {
		http.Error(w, "short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func (h *URLHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	data := h.Repo.GetTopDomains(3)
	json.NewEncoder(w).Encode(data)
}
