package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brij-812/url-shortener/internal/repository"
)

// fakeUserID we’ll inject into request context to simulate JWT middleware
const fakeUserID = 42

func withUserContext(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), "user_id", fakeUserID)
	return req.WithContext(ctx)
}

func TestShortenAndMetrics(t *testing.T) {
	repo := repository.NewMemoryRepo()
	h := NewURLHandler(repo)

	// ---- test shorten ----
	body := []byte(`{"url":"https://a.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	req = withUserContext(req) // simulate authenticated user

	w := httptest.NewRecorder()
	h.ShortenURL(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}

	// ---- test metrics ----
	req2 := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req2 = withUserContext(req2)
	w2 := httptest.NewRecorder()

	h.GetMetrics(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("unexpected metrics status %d", w2.Code)
	}

	var data map[string]int
	if err := json.Unmarshal(w2.Body.Bytes(), &data); err != nil {
		t.Fatalf("failed to unmarshal metrics: %v", err)
	}

	if data["a.com"] != 1 {
		t.Fatalf("expected a.com count 1, got %v", data)
	}

	// ---- test GetAllUserURLs ----
	req3 := httptest.NewRequest(http.MethodGet, "/all", nil)
	req3 = withUserContext(req3)
	w3 := httptest.NewRecorder()

	h.GetAllUserURLs(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("unexpected /all status %d", w3.Code)
	}

	var urls []map[string]string
	if err := json.Unmarshal(w3.Body.Bytes(), &urls); err != nil {
		t.Fatalf("failed to unmarshal urls: %v", err)
	}
	if len(urls) == 0 {
		t.Fatalf("expected at least one URL entry, got empty list")
	}
}
