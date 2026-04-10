package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User credentials for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type AuthError struct {
	Error string `json:"error"`
}

// JWTClaims represents the JWT claims
type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Secret key for JWT (in production, use environment variable)
const jwtSecret = "your-super-secret-jwt-key-change-in-production"

// Sample user credentials (in production, use database)
var validUsers = map[string]string{
	"shyvanshu.mehra@triconinfotech.com": hashPassword("Tricon@99"),
	"admin@example.com":                  hashPassword("admin123"),
	"user@example.com":                   hashPassword("user123"),
	"test@example.com":                   hashPassword("test123"),
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// AuthHandler handles authentication endpoints
type AuthHandler struct{}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login handles user login and JWT token generation
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthError{Error: "Invalid request format"})
		return
	}

	// Validate credentials
	if req.Email == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthError{Error: "Email and password are required"})
		return
	}

	// Check credentials
	storedHash, exists := validUsers[req.Email]
	passwordHash := hashPassword(req.Password)

	if !exists || storedHash != passwordHash {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthError{Error: "Invalid email or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthError{Error: "Failed to generate token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Token:   tokenString,
		Email:   req.Email,
		Message: "Login successful",
	})
}

// Verify token and return claims
func (ah *AuthHandler) VerifyToken(tokenString string) (*JWTClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// Middleware to check JWT token
func (ah *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: AuthMiddleware - checking authorization")
		authHeader := r.Header.Get("Authorization")
		log.Printf("DEBUG: AuthMiddleware - authorization header: %s", authHeader)

		if authHeader == "" {
			log.Println("DEBUG: AuthMiddleware - missing authorization header")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthError{Error: "Missing authorization header"})
			return
		}

		claims, err := ah.VerifyToken(authHeader)
		if err != nil {
			log.Printf("DEBUG: AuthMiddleware - token verification error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthError{Error: "Invalid or expired token"})
			return
		}

		log.Printf("DEBUG: AuthMiddleware - token valid for user: %s", claims.Email)
		// Store email in context for later use
		r.Header.Set("X-User-Email", claims.Email)
		next(w, r)
	}
}
