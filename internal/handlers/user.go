package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/brij-812/url-shortener/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("supersecretkey")

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// POST /signup
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec(`INSERT INTO users (email, password_hash, created_at) VALUES ($1, $2, $3)`,
		req.Email, string(hashed), time.Now())
	if err != nil {
		http.Error(w, "email already exists or db error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// POST /login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var storedHash string
	var userID int
	err := h.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email = $1`, req.Email).Scan(&userID, &storedHash)
	if err == sql.ErrNoRows {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.AuthResponse{Token: tokenString})
}
