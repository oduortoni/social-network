package service

import (
	"time"
	"net/http"

	"fmt"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// AuthService handles the business logic for authentication.
type AuthService struct {
	AuthStore *store.AuthStore
}

const (
	EXPIRED_SESSION = "Session expired"
	INVALID_PASSWORD = "Invalid password"
	INVALID_EMAIL    = "Invalid email"
)

// NewAuthService creates a new AuthService.
func NewAuthService(as *store.AuthStore) *AuthService {
	return &AuthService{AuthStore: as}
}

// AuthenticateUser authenticates a user by email and password.
func (s *AuthService) AuthenticateUser(email, password string) (*models.User, string, error) {
	user, err := s.AuthStore.GetUserByEmail(email)
	if err != nil {
		return nil, INVALID_EMAIL, err
	}

	fmt.Println("User found:", user)

	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	if err := passwordManager.ComparePassword(user.Password, password); err != nil {
		fmt.Println("Password mismatch for user:", user.Email)
		return nil, INVALID_PASSWORD, err
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := s.AuthStore.CreateSession(user.ID, sessionID, expiresAt); err != nil {
		fmt.Println("Failed to create session:", err)
		return nil, EXPIRED_SESSION, err
	}

	return user, sessionID, nil
}

func (s *AuthService) DeleteSession(sessionID string) (int, error) {
	if err := s.AuthStore.DeleteSession(sessionID); err != nil {
		return 0, err
	}
	return http.StatusOK, nil
}

func (s *AuthService) GetUserIDBySession(sessionID string) (int, error) {
	userID, err := s.AuthStore.GetUserIDBySession(sessionID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
