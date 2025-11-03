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
	Repo repository.Repository
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

// 🔹 Shorten a new URL (Protected)
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
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	code, exists := h.Repo.GetCode(req.URL)
	if !exists {
		code = utils.GenerateShortCode(req.URL)
		h.Repo.Save(req.URL, code, userID)
	} else {
		h.Repo.IncrementDomainCount(req.URL, userID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ShortenResponse{
		ShortURL: "http://localhost:8080/" + code,
	})
}

// 🔹 Redirect (Public)
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.Error(w, "short code missing", http.StatusBadRequest)
		return
	}

	longURL, exists := h.Repo.GetURL(shortCode)
	if !exists {
		http.Error(w, "short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

// 🔹 Metrics (Protected, per user)
func (h *URLHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	data := h.Repo.GetTopDomains(userID, 3)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// 🔹 Get all URLs created by the current user (Protected)
func (h *URLHandler) GetAllUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	urls := h.Repo.GetAllURLsByUser(userID)
	w.Header().Set("Content-Type", "application/json")

	if len(urls) == 0 {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	json.NewEncoder(w).Encode(urls)
}
