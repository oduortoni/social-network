package authentication

import (
	"net/http"
	"time"
)

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Delete session from the DB

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Value = ""
	cookie.Path = "/"

	// Update the cookie in the response
	http.SetCookie(w, cookie)

	respondJSON(w, http.StatusOK, Response{Message: "Logout successful"})
}
