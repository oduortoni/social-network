package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockAuthService is a mock implementation of the AuthService for testing.
type MockAuthService struct {
	AuthenticateUserFunc func(email, password string) (*models.User, string, error)
}

func (s *MockAuthService) AuthenticateUser(email, password string) (*models.User, string, error) {
	return s.AuthenticateUserFunc(email, password)
}

func TestLogin(t *testing.T) {
	// Create a new mock auth service
	mockAuthService := &MockAuthService{
		AuthenticateUserFunc: func(email, password string) (*models.User, string, error) {
			return &models.User{ID: 1, Email: "test@test.com"}, "session123", nil
		},
	}

	// Create a new auth handler with the mock service
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Create a new request
	loginReq := handlers.LoginRequest{Email: "test@test.com", Password: "password"}
	body, _ := json.Marshal(loginReq)
	req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder
	rr := httptest.NewRecorder()

	// Call the handler
	authHandler.Login(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var user models.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if user.Email != "test@test.com" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			user.Email, "test@test.com")
	}
}
