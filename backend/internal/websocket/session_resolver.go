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

/*
*  Gets the user id from the request ans uses it to get the user's nickname and avatar from the database.
*  Returns the user id, nickname, and avatar.
*/
func (r *DBSessionResolver) GetUserFromRequest(req *http.Request) (int64, string, string, error) {
	var sessionID string

	// First try to get session ID from cookie
	cookie, err := req.Cookie("session_id")
	if err != nil {
		// If cookie fails, try query parameter as fallback
		sessionID = req.URL.Query().Get("session_id")
		if sessionID == "" {
			return 0, "anonymous", "", err
		}
	} else {
		sessionID = cookie.Value
	}

	var userID int64
	var nickname string
	var avatar string
	err = r.DB.QueryRow(`
		SELECT sessions.user_id, users.nickname, users.avatar FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE sessions.id = ?
	`, sessionID).Scan(&userID, &nickname, &avatar)
	if err != nil {
		return 0, "anonymous", "", err
	}
	return userID, nickname, avatar, nil
}
