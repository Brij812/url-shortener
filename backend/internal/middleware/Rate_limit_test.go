package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brij-812/HyperLinkOS/internal/cache"
)

func TestRateLimitMiddleware(t *testing.T) {
	// make sure Redis is running locally
	cache.InitRedis("localhost:6379", "", 0)

	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// helper to make one request
	doRequest := func(ip string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("GET", "/shorten", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w
	}

	// hit the same IP repeatedly up to the limit
	ip := "127.0.0.1:10000"
	for i := 1; i <= userLimit; i++ {
		resp := doRequest(ip)
		if resp.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for request %d, got %d", i, resp.Code)
		}
	}

	// next request from same IP should trigger rate limit
	resp := doRequest(ip)
	if resp.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", resp.Code)
	}

	// verify headers exist
	if resp.Header().Get("X-RateLimit-Limit") == "" {
		t.Errorf("missing X-RateLimit-Limit header")
	}
	if resp.Header().Get("X-RateLimit-Remaining") == "" {
		t.Errorf("missing X-RateLimit-Remaining header")
	}
	if resp.Header().Get("X-RateLimit-Reset") == "" {
		t.Errorf("missing X-RateLimit-Reset header")
	}

	t.Logf("✅ Rate limit headers verified successfully")
}

func TestRateLimitReset(t *testing.T) {
	cache.InitRedis("localhost:6379", "", 0)

	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ip := "127.0.0.1:20000"

	// trigger the limit
	for i := 0; i < userLimit; i++ {
		req := httptest.NewRequest("GET", "/shorten", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	// one more before reset should fail
	req := httptest.NewRequest("GET", "/shorten", nil)
	req.RemoteAddr = ip
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 before reset, got %d", w.Code)
	}

	// wait for window to reset
	t.Log("waiting for rate limit window to reset...")
	time.Sleep(time.Duration(windowSecs+2) * time.Second)

	req = httptest.NewRequest("GET", "/shorten", nil)
	req.RemoteAddr = ip
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 after window reset, got %d", w.Code)
	}
}
