package authentication

import (
	"database/sql"
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

type SigninResponse struct {
	Message string `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req User
	var serverresponse SigninResponse = SigninResponse{
		Message: "Login successful",
		Data:    map[string]any{},
	}
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
	userInDB, err := GetUserByEmail(req.Email, db)
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
	err = DeleteUserSessions(userInDB.ID, db)
	if err != nil {
		serverresponse.Message = "Failed to delete previous sessions"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Create session
	sessionID := uuid.New().String()
	expirey := time.Now().Add(24 * time.Hour)

	// Set cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expirey,
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	}
	if err := StoreSession(userInDB.ID, sessionID, expirey, db); err != nil {
		serverresponse.Message = "Failed to create session"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	http.SetCookie(w, &cookie)

	// Success response
	serverresponse.Message = "Login successful"
	serverresponse.Data = map[string]any{
		"redirect": "/dashboard",
	}
	respondJSON(w, statusCode, serverresponse)
}

func GetUserByEmail(email string, db *sql.DB) (User, error) {
	var users User

	err := db.QueryRow("SELECT id, email, password FROM Users WHERE email = ?", email).Scan(&users.ID, &users.Email, &users.Password)
	if err != nil {
		return User{}, err // Return an error if the query fails
	}

	return users, nil
}
