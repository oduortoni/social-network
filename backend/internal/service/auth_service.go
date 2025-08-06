package service

import (
	"fmt"
	"net/http"
	"regexp"
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

const (
	EXPIRED_SESSION  = "Session expired"
	INVALID_PASSWORD = "Invalid password"
	INVALID_EMAIL    = "User does not exist"
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

// CreateUser creates a new user with validation and password hashing
func (s *AuthService) CreateUser(user *models.User) (*models.User, error) {
	// Check if user already exists
	exists, err := s.AuthStore.UserExists(user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Hash the password
	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	hashedPassword, err := passwordManager.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Set created time
	user.CreatedAt = time.Now()

	// Create user in database
	userID, err := s.AuthStore.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = userID
	return user, nil
}

// validateEmail validates email format using regex
func (s *AuthService) ValidateEmail(email string) (bool, error) {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re, err := regexp.Compile(emailPattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(email), nil
}

func (s *AuthService) UserExists(email string) (bool, error) {
	return s.AuthStore.UserExists(email)
}

func (s *AuthService) UserNewEditEmailExist(email string, userid int64) (bool, error) {
	return s.AuthStore.NewEditEmailExist(email, userid)
}


func (s *AuthService)EditUserProfile(user *models.User, userid int64) (error){
	  return  s.AuthStore.EditProfile(user,userid)
}