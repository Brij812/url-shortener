package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/brij-812/url-shortener/internal/cache"
)

// window length and limits
const (
	windowSecs = 60
	userLimit  = 10
	ipLimit    = 30
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		window := now / windowSecs

		userID := r.Context().Value("user_id")
		var keyBase string
		if userID != nil {
			keyBase = fmt.Sprintf("rate:user:%v", userID)
		} else {
			keyBase = fmt.Sprintf("rate:ip:%s", clientIP(r))
		}

		currKey := fmt.Sprintf("%s:%d", keyBase, window)
		prevKey := fmt.Sprintf("%s:%d", keyBase, window-1)

		pipe := cache.Client().TxPipeline()
		currCount := pipe.Incr(r.Context(), currKey)
		pipe.Expire(r.Context(), currKey, time.Duration(windowSecs*2)*time.Second)
		_, _ = pipe.Exec(r.Context())

		prevVal, _ := cache.Client().Get(r.Context(), prevKey).Int64()
		currVal := currCount.Val()

		elapsed := float64(now%windowSecs) / float64(windowSecs)
		blended := float64(prevVal)*(1.0-elapsed) + float64(currVal)

		limit := userLimit
		if userID == nil {
			limit = ipLimit
		}

		if blended > float64(limit) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded. Try again later."))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// clientIP extracts the real client IP.
func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
