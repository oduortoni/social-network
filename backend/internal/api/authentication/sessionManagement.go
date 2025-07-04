package authentication

import (
	"database/sql"
	"net/http"
	"time"
)

// CheckSessionHandler verifies if a session exists and is valid in the database
func CheckSessionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, Response{Message: "No session found"})
		return
	}

	var expiresAt time.Time
	err = db.QueryRow("SELECT expires_at FROM Sessions WHERE id = ?", cookie.Value).Scan(&expiresAt)
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
		_, err := db.Exec("DELETE FROM Sessions WHERE id = ?", cookie.Value)
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
