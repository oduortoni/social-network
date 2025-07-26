package websocket

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type DBSessionResolver struct {
	DB *sql.DB
}

func NewDBSessionResolver(db *sql.DB) *DBSessionResolver {
	return &DBSessionResolver{DB: db}
}

func (r *DBSessionResolver) GetUserIDFromRequest(req *http.Request) (int64, string, error) {
	fmt.Println("Getting user ID from request...")
	var sessionID string

	// First try to get session ID from cookie
	cookie, err := req.Cookie("session_id")
	if err != nil {
		// If cookie fails, try query parameter as fallback
		sessionID = req.URL.Query().Get("session_id")
		if sessionID == "" {
			return 0, "anonymous", err
		}
	} else {
		sessionID = cookie.Value
	}

	fmt.Println("Session ID:", sessionID)

	var userID int64
	var nickname string
	err = r.DB.QueryRow(`
		SELECT sessions.user_id, users.nickname FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE sessions.id = ?
	`, sessionID).Scan(&userID, &nickname)
	if err != nil {
		log.Fatalf("Failed to query user ID from session: %v", err)
		return 0, "anonymous", err
	}
	fmt.Println("Nickname:", nickname)
	return userID, nickname, nil
}
