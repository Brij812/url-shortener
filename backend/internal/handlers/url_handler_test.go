package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brij-812/HyperLinkOS/internal/repository"
)

func TestShortenAndMetrics(t *testing.T) {
	repo := repository.NewMemoryRepo()
	h := NewURLHandler(repo)

	body := []byte(`{"url":"https://a.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ShortenURL(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w2 := httptest.NewRecorder()
	h.GetMetrics(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("unexpected metrics status %d", w2.Code)
	}

	var data map[string]int
	json.Unmarshal(w2.Body.Bytes(), &data)
	if data["a.com"] != 1 {
		t.Fatalf("expected a.com count 1, got %v", data)
	}
}
