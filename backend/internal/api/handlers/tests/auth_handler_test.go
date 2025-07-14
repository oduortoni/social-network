package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockAuthService is a mock implementation of the AuthService for testing.
type MockAuthService struct {
	AuthenticateUserFunc func(email, password string) (*models.User, string, error)
	DeleteSessionFunc    func(sessionID string) (int, error)
	GetUserIDBySessionFunc func(sessionID string) (int, error)
}

func (s *MockAuthService) AuthenticateUser(email, password string) (*models.User, string, error) {
	return s.AuthenticateUserFunc(email, password)
}

func (s *MockAuthService) DeleteSession(sessionID string) (int, error) {
	if s.DeleteSessionFunc != nil {
		return s.DeleteSessionFunc(sessionID)
	}
	return 0, nil
}

func (s *MockAuthService) GetUserIDBySession(sessionID string) (int, error) {
	return s.GetUserIDBySessionFunc(sessionID)
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
	req.Header.Set("Content-Type", "application/json")

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
	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Logged in successfully"
	if resp.Message != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Message, expectedMessage)
	}

	// Check for session cookie
	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			found = true
			if cookie.Value != "session123" {
				t.Errorf("cookie session_id has wrong value: got %q, want %q", cookie.Value, "session123")
			}
			if cookie.Path != "/" {
				t.Errorf("cookie session_id has wrong path: got %q, want %q", cookie.Path, "/")
			}
			if !cookie.HttpOnly {
				t.Error("cookie session_id is not http-only")
			}
			if cookie.SameSite != http.SameSiteLaxMode {
				t.Errorf("cookie session_id has wrong same-site policy: got %v, want %v", cookie.SameSite, http.SameSiteLaxMode)
			}
			if time.Until(cookie.Expires) > 7*24*time.Hour {
				t.Error("cookie session_id has wrong expiration")
			}
		}
	}

	if !found {
		t.Error("session_id cookie not set")
	}
}
