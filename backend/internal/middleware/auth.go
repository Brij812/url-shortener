package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWTSecret initializes the global JWT secret
func InitJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// JWTAuth validates the JWT (from cookie or Authorization header)
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// 1️⃣ Prefer cookie (browser clients)
		if cookie, err := r.Cookie("hl_jwt"); err == nil {
			tokenString = cookie.Value
		} else {
			// 2️⃣ Fallback to Authorization header for API tools
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			http.Error(w, "missing auth token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		// ✅ Check token expiry
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
		}

		// ✅ Extract user_id
		var userID int
		switch v := claims["user_id"].(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		default:
			http.Error(w, "invalid user_id in token", http.StatusUnauthorized)
			return
		}

		log.Printf("✅ Authenticated request by user_id=%d", userID)

		// ✅ Inject into context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
