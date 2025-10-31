package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/brij-812/url-shortener/internal/cache"
)

const (
	userLimit  = 10
	ipLimit    = 30
	windowSecs = 60
)

// RateLimit middleware controls request frequency per user/IP.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key string

		// Extract user ID if logged in
		userID := r.Context().Value("user_id")
		if userID != nil {
			key = fmt.Sprintf("rate:user:%v", userID)
		} else {
			ip := clientIP(r)
			key = fmt.Sprintf("rate:ip:%s", ip)
		}

		count, _ := cache.Client().Incr(r.Context(), key).Result()

		if count == 1 {
			// first hit — set expiration window
			cache.Client().Expire(r.Context(), key, time.Duration(windowSecs)*time.Second)
		}

		limit := userLimit
		if userID == nil {
			limit = ipLimit
		}

		if count > int64(limit) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded. Try again later."))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// clientIP extracts IP from request.
func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
