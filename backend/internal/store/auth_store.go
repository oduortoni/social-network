package store

import (
	"database/sql"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// AuthStore handles database operations for authentication.
type AuthStore struct {
	DB *sql.DB
}

// NewAuthStore creates a new AuthStore.
func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{DB: db}
}

// GetUserByEmail fetches a user by email
func (s *AuthStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.DB.QueryRow(
		"SELECT id, email, password FROM Users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// CreateSession creates a new session record and returns the session ID
/*
func (s *UserStore) CreateSession(userID int) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days expiry

	_, err := s.db.Exec(
		"INSERT INTO Sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
*/
// CreateSession creates a new session for a user.
func (s *AuthStore) CreateSession(userID int64, sessionID string, expiresAt time.Time) error {
	_, err := s.DB.Exec("INSERT INTO Sessions (user_id, id, created_at, expires_at) VALUES (?, ?, ?, ?)", userID, sessionID, time.Now(), expiresAt)
	return err
}

// -----------------------------------------------

// GetUserIDBySession returns the user ID associated with a valid session ID
func (s *AuthStore) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := s.DB.QueryRow(
		"SELECT user_id, expires_at FROM Sessions WHERE id = ?",
		sessionID,
	).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if expiresAt.Before(time.Now()) {
		return 0, sql.ErrNoRows // session expired
	}

	return userID, nil
}

// DeleteSession deletes a session by ID
func (s *AuthStore) DeleteSession(sessionID string) error {
	_, err := s.DB.Exec("DELETE FROM Sessions WHERE id = ?", sessionID)
	return err
}
