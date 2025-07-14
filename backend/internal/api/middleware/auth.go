package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"github.com/tajjjjr/social-network/backend/internal/api/authentication"

	"fmt"
)

// AuthMiddleware retrieves the user ID from the session cookie and adds it to the request context.
func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			fmt.Println("AuthMiddleware: Cookie retrieved:", cookie)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := authentication.GetUserIDFromSession(cookie.Value, db)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
