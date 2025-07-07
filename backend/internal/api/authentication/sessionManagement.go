package authentication

import (
	"database/sql"
	"net/http"
	"time"
)

// GetUserIDFromSession retrieves the user ID associated with a given session ID.
func GetUserIDFromSession(sessionID string, db *sql.DB) (int64, error) {
	var userID int64
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// CheckSessionHandler verifies if a session exists and is valid in the database
func CheckSessionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, Response{Message: "No session found"})
		return
	}

	expiresAt, err := GetSessionsExpiretime(cookie, db)
	if err != nil {
		if err == sql.ErrNoRows {
			respondJSON(w, http.StatusUnauthorized, Response{Message: "Invalid session"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, Response{Message: "Database error"})
		return
	}

	// Check if session has expired
	if time.Now().After(expiresAt) {
		// Delete expired session
		err = DeleteSession(cookie, db)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, Response{Message: "Failed to delete expired session"})
			return
		}
		respondJSON(w, http.StatusUnauthorized, Response{Message: "Session expired"})
		return
	}

	respondJSON(w, http.StatusOK, Response{Message: "Valid session"})
}

func DeleteUserSessions(id int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Sessions WHERE user_id = ?", id)

	return err
}

func StoreSession(id int, sessionID string, expiray time.Time, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO Sessions (user_id, id, created_at,expires_at) VALUES (?, ?, ?,?)", id, sessionID, time.Now(), expiray)
	return err
}

func DeleteSession(cookie *http.Cookie, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Sessions WHERE id = ?", cookie.Value)

	return err
}

func GetSessionsExpiretime(cookie *http.Cookie, db *sql.DB) (time.Time, error) {
	var expiresAt time.Time
	err := db.QueryRow("SELECT expires_at FROM Sessions WHERE id = ?", cookie.Value).Scan(&expiresAt)
	return expiresAt, err
}
