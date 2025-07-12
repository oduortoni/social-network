package controllers

import (
	"database/sql"
	"net/http"
	auth "github.com/tajjjjr/social-network/backend/internal/api/authentication"
)

func Me(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	auth.CheckSessionHandler(w, r, db)
}
