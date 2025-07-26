package websocket

import (
	"database/sql"
	"net/http"
)

type DBSessionResolver struct {
	DB *sql.DB
}

func NewDBSessionResolver(db *sql.DB) *DBSessionResolver {
	return &DBSessionResolver{DB: db}
}

func (r *DBSessionResolver) GetUserIDFromRequest(req *http.Request) (int64, string, error) {
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

	var userID int64
	var nickname string
	err = r.DB.QueryRow(`
		SELECT sessions.user_id, users.nickname FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE sessions.id = ?
	`, sessionID).Scan(&userID, &nickname)
	if err != nil {
		return 0, "anonymous", err
	}
	return userID, nickname, nil
}
