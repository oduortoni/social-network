package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// AuthService handles the business logic for authentication.
type AuthService struct {
	AuthStore *store.AuthStore
}

// NewAuthService creates a new AuthService.
func NewAuthService(as *store.AuthStore) *AuthService {
	return &AuthService{AuthStore: as}
}

// AuthenticateUser authenticates a user by email and password.
func (s *AuthService) AuthenticateUser(email, password string) (*models.User, string, error) {
	user, err := s.AuthStore.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}

	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	if err := passwordManager.ComparePassword(user.Password, password); err != nil {
		return nil, "", nil // Invalid password
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := s.AuthStore.CreateSession(user.ID, sessionID, expiresAt); err != nil {
		return nil, "", err
	}

	return user, sessionID, nil
}
