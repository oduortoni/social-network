package authentication

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req User
	var serverresponse Response
	statusCode := http.StatusOK

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		serverresponse.Message = "Failed to read request"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Parse JSON
	if err := json.Unmarshal(body, &req); err != nil {
		serverresponse.Message = "Invalid JSON"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Get user from database
	userInDB, err := GetUserByEmail(req.Email)
	if err != nil {
		serverresponse.Message = "User not found"
		statusCode = http.StatusUnauthorized
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Compare password
	if !CompareHashpassword(userInDB.Password, req.Password) {
		serverresponse.Message = "Invalid credentials"
		statusCode = http.StatusUnauthorized
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Clear previous sessions
	DeleteUserSessions(userInDB.ID)

	// Create session
	sessionID := uuid.New().String()
	if err := StoreSession(userInDB.ID, sessionID); err != nil {
		serverresponse.Message = "Failed to create session"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Set cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	// Success response
	serverresponse.Message = "Login successful"
	statusCode = http.StatusOK
	respondJSON(w, statusCode, serverresponse)
}

func DeleteUserSessions(id int) {
}

func StoreSession(id int, sessionID string) error {
	return nil
}

func GetUserByEmail(email string) (User, error) {
	return User{}, nil
}
