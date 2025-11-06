package middleware

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/brij-812/HyperLinkOS/internal/cache"
)

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

		// increment current counter
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

		remaining := int64(math.Max(0, float64(limit)-blended))
		resetIn := windowSecs - int(now%windowSecs)

		// 🔹 Standard Rate-Limit headers (like GitHub)
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetIn))

		if blended > float64(limit) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf("Rate limit exceeded. Try again in %d seconds.", resetIn)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
