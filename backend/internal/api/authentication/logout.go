package authentication

import (
	"database/sql"
	"net/http"
	"time"
)

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, Response{Message: "No session found"})
		return
	}

	err = DeleteSession(cookie, db)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Message: "Failed to delete session"})
		return
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Value = ""
	cookie.Path = "/"

	http.SetCookie(w, cookie)

	respondJSON(w, http.StatusOK, Response{Message: "Logout successful"})
}

