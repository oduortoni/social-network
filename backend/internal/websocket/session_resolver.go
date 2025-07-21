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

func (r *DBSessionResolver) GetUserIDFromRequest(req *http.Request) (int64, error) {
	cookie, err := req.Cookie("session_id")
	if err != nil {
		return 0, err
	}

	var userID int64
	err = r.DB.QueryRow(`
		SELECT user_id FROM sessions WHERE id = ?
	`, cookie.Value).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
