package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/utils"
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

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	var err error

	// Check Content-Type header to determine how to parse the request
	contentType := r.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	if strings.Contains(contentType, "application/json") {
		// Parse JSON body
		err = json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid JSON request body"})
			return
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		err = r.ParseForm()
		if err != nil {
			models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid form data"})
			return
		}

		// Extract form values
		creds.Email = r.FormValue("email")
		creds.Password = r.FormValue("password")
	} else {
		// Default to trying JSON first, then form data
		err = json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			// If JSON fails, try to parse as form data
			if parseErr := r.ParseForm(); parseErr != nil {
				models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid request body format"})
				return
			}
			creds.Email = r.FormValue("email")
			creds.Password = r.FormValue("password")
		}
	}

	fmt.Println("Login credentials:", creds)

	authUser, sessionID, _ := auth.AuthService.AuthenticateUser(creds.Email, creds.Password)
	if authUser == nil {
		if sessionID == service.EXPIRED_SESSION {
			models.RespondJSON(w, http.StatusInternalServerError, models.Response{Message: "Failed to create session"})
		} else if sessionID == service.INVALID_PASSWORD {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "Invalid password"})
		} else if sessionID == service.INVALID_EMAIL {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "User not found"})
		}
		return
	}

	fmt.Println("Authenticated user:", authUser)

	// only used for UI checks to avoid flashing protected routes
	fmt.Println("Setting login confirmation cookie:")
	http.SetCookie(w, &http.Cookie{
		Name:  "logged_in",
		Value: "true",
		Path:  "/",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Secure: true,
	})

	// used to actually authenticate users
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	}
	fmt.Println("Setting session cookie:", sessionCookie)
	http.SetCookie(w, sessionCookie)
	models.RespondJSON(w, http.StatusOK, models.Response{Message: "Logged in successfully"})
}

// LogoutHandler deletes session and clears cookie
func (auth *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_, _ = auth.AuthService.DeleteSession(cookie.Value)

		// Clear cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false, // set true in production
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "logged_in",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
			Secure:   true,
		})

	}

	models.RespondJSON(w, http.StatusOK, models.Response{Message: "Logged out successfully"})
}

// AuthMiddleware verifies session cookie, loads user ID into context
func (auth *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "Authentication required"})
			return
		}

		userID, err := auth.AuthService.GetUserIDBySession(cookie.Value)
		if err != nil {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "Invalid or expired session"})
			return
		}

		// Add userID to request context for downstream handlers
		ctx := context.WithValue(r.Context(), utils.User_id, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
