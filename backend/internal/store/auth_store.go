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

// GetUserByEmail retrieves a user by their email address.
func (s *AuthStore) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow("SELECT id, email, password FROM Users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateSession creates a new session for a user.
func (s *AuthStore) CreateSession(userID int64, sessionID string, expiresAt time.Time) error {
	_, err := s.DB.Exec("INSERT INTO Sessions (user_id, id, created_at, expires_at) VALUES (?, ?, ?, ?)", userID, sessionID, time.Now(), expiresAt)
	return err
}
