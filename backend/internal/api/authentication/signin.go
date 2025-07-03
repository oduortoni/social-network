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

func SigninHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if err := StoreSession(userInDB.ID, sessionID, db); err != nil {
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

func DeleteUserSessions(id int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Sessions WHERE user_id = ?", id) 

	return err
}

func StoreSession(id int, sessionID string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO Sessions (user_id, session_id, created_at) VALUES (?, ?, ?)", id, sessionID, time.Now())
	return err
}

func GetUserByEmail(email string, db *sql.DB) (User, error) {
	var users User

	err := db.QueryRow("SELECT id, email, password FROM Users WHERE email = ?", email).Scan(&users.ID, &users.Email, &users.Password)
	if err != nil {
		return User{}, err // Return an error if the query fails
	}

	return users, nil
}
