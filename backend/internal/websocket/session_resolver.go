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
	cookie, err := req.Cookie("session_id")
	if err != nil {
		return 0, "anonymous", err
	}

	var userID int64
	var nickname string
	err = r.DB.QueryRow(`
		SELECT user_id, users.nickname FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE id = ?
	`, cookie.Value).Scan(&userID, &nickname)
	if err != nil {
		return 0, "anonymous", err
	}
	return userID, nickname, nil
}
