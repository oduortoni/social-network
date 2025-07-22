package tests

import (
	"encoding/json"
	"errors"
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
	AuthenticateUserFunc   func(email, password string) (*models.User, string, error)
	DeleteSessionFunc      func(sessionID string) (int, error)
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

func TestLogin_SessionFixation_NotPrevented(t *testing.T) {
	mockAuthService := &MockAuthService{
		AuthenticateUserFunc: func(email, password string) (*models.User, string, error) {
			return &models.User{ID: 1, Email: "test@test.com"}, "new-session-id", nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	loginReq := handlers.LoginRequest{Email: "test@test.com", Password: "password"}
	body, _ := json.Marshal(loginReq)
	req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	preSessionID := "fixed-session-id"
	req.AddCookie(&http.Cookie{Name: "session_id", Value: preSessionID})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	authHandler.Login(rr, req)

	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			found = true
			if cookie.Value == preSessionID {
				t.Error("Session fixation detected: session_id was not regenerated on login")
			}
		}
	}
	if !found {
		t.Error("session_id cookie not set after login")
	}
}

// Test login with incorrect credentials
func TestLogin_IncorrectCredentials(t *testing.T) {
	testCases := []struct {
		name     string
		email    string
		password string
		mockFunc func(email, password string) (*models.User, string, error)
	}{
		{
			name:     "Invalid email",
			email:    "nonexistent@test.com",
			password: "password",
			mockFunc: func(email, password string) (*models.User, string, error) {
				return nil, "User does not exist", errors.New("user not found")
			},
		},
		{
			name:     "Invalid password",
			email:    "test@test.com",
			password: "wrongpassword",
			mockFunc: func(email, password string) (*models.User, string, error) {
				return nil, "Invalid password", errors.New("password mismatch")
			},
		},
		{
			name:     "Empty email",
			email:    "",
			password: "password",
			mockFunc: func(email, password string) (*models.User, string, error) {
				return nil, "User does not exist", errors.New("empty email")
			},
		},
		{
			name:     "Empty password",
			email:    "test@test.com",
			password: "",
			mockFunc: func(email, password string) (*models.User, string, error) {
				return nil, "Invalid password", errors.New("empty password")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{
				AuthenticateUserFunc: tc.mockFunc,
			}
			authHandler := handlers.NewAuthHandler(mockAuthService)

			loginReq := handlers.LoginRequest{Email: tc.email, Password: tc.password}
			body, _ := json.Marshal(loginReq)
			req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			authHandler.Login(rr, req)

			// Should return unauthorized for incorrect credentials
			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusUnauthorized)
			}

			// Should not set session cookie for failed login
			cookies := rr.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "session_id" && cookie.Value != "" {
					t.Error("session_id cookie should not be set for failed login")
				}
			}
		})
	}
}

// Test login with form data instead of JSON
func TestLogin_FormData(t *testing.T) {
	mockAuthService := &MockAuthService{
		AuthenticateUserFunc: func(email, password string) (*models.User, string, error) {
			if email == "test@test.com" && password == "password" {
				return &models.User{ID: 1, Email: "test@test.com"}, "session123", nil
			}
			return nil, "Invalid credentials", errors.New("authentication failed")
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Test with form data
	formData := "email=test@test.com&password=password"
	req, err := http.NewRequest("POST", "/login", strings.NewReader(formData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	authHandler.Login(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check for session cookie
	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" && cookie.Value == "session123" {
			found = true
		}
	}
	if !found {
		t.Error("session_id cookie not set correctly")
	}
}

// Test session creation properties
func TestLogin_SessionCookieProperties(t *testing.T) {
	mockAuthService := &MockAuthService{
		AuthenticateUserFunc: func(email, password string) (*models.User, string, error) {
			return &models.User{ID: 1, Email: "test@test.com"}, "session123", nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	loginReq := handlers.LoginRequest{Email: "test@test.com", Password: "password"}
	body, _ := json.Marshal(loginReq)
	req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	authHandler.Login(rr, req)

	// Check session cookie properties
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Fatal("session_id cookie not found")
	}

	// Verify cookie properties
	if sessionCookie.Path != "/" {
		t.Errorf("cookie path should be '/', got %q", sessionCookie.Path)
	}
	if !sessionCookie.HttpOnly {
		t.Error("cookie should be HttpOnly")
	}
	if sessionCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("cookie SameSite should be Lax, got %v", sessionCookie.SameSite)
	}
	if sessionCookie.Expires.IsZero() {
		t.Error("cookie should have expiration time set")
	}
	// Check that expiration is approximately 7 days from now
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	if sessionCookie.Expires.Before(expectedExpiry.Add(-time.Minute)) ||
		sessionCookie.Expires.After(expectedExpiry.Add(time.Minute)) {
		t.Errorf("cookie expiration should be ~7 days from now, got %v", sessionCookie.Expires)
	}
}
