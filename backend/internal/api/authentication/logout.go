package authentication

import (
	"database/sql"
	"net/http"
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

	// Create a new cookie to clear the client's session cookie.
	// Setting MaxAge to -1 is the most reliable way to tell the browser to delete it.
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	respondJSON(w, http.StatusOK, Response{Message: "Logout successful, session invalidated"})
}
