package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/service"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	AuthService service.AuthServiceInterface
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(as service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{AuthService: as}
}

// LoginRequest represents the request body for a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, sessionID, err := h.AuthService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
