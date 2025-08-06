package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// Test session persistence by accessing protected routes
func TestSessionPersistence_ValidSession(t *testing.T) {
	mockAuthService := &MockAuthService{
		GetUserIDBySessionFunc: func(sessionID string) (int, error) {
			if sessionID == "valid-session-123" {
				return 1, nil
			}
			return 0, errors.New("invalid session")
		},
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Create a request with a valid session cookie
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "valid-session-123",
	})

	rr := httptest.NewRecorder()

	// Create a dummy protected handler
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
	})

	// Wrap with auth middleware
	authMiddleware := authHandler.AuthMiddleware(protectedHandler)
	authMiddleware.ServeHTTP(rr, req)

	// Should allow access
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Access granted" {
		t.Errorf("expected 'Access granted', got %q", resp.Message)
	}
}

func TestSessionPersistence_InvalidSession(t *testing.T) {
	mockAuthService := &MockAuthService{
		GetUserIDBySessionFunc: func(sessionID string) (int, error) {
			return 0, errors.New("invalid session")
		},
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Create a request with an invalid session cookie
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "invalid-session-456",
	})

	rr := httptest.NewRecorder()

	// Create a dummy protected handler
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
	})

	// Wrap with auth middleware
	authMiddleware := authHandler.AuthMiddleware(protectedHandler)
	authMiddleware.ServeHTTP(rr, req)

	// Should deny access
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Invalid or expired session" {
		t.Errorf("expected 'Invalid or expired session', got %q", resp.Message)
	}
}

func TestSessionPersistence_NoSession(t *testing.T) {
	mockAuthService := &MockAuthService{
		GetUserIDBySessionFunc: func(sessionID string) (int, error) {
			return 0, errors.New("no session")
		},
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Create a request without session cookie
	req := httptest.NewRequest("GET", "/protected", nil)

	rr := httptest.NewRecorder()

	// Create a dummy protected handler
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
	})

	// Wrap with auth middleware
	authMiddleware := authHandler.AuthMiddleware(protectedHandler)
	authMiddleware.ServeHTTP(rr, req)

	// Should deny access
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Authentication required" {
		t.Errorf("expected 'Authentication required', got %q", resp.Message)
	}
}

// Test logout functionality
func TestLogout_ValidSession(t *testing.T) {
	mockAuthService := &MockAuthService{
		DeleteSessionFunc: func(sessionID string) (int, error) {
			if sessionID == "valid-session-123" {
				return http.StatusOK, nil
			}
			return 0, errors.New("session not found")
		},
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Create a request with a valid session cookie
	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "valid-session-123",
	})

	rr := httptest.NewRecorder()
	authHandler.LogoutHandler(rr, req)

	// Should return success
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that cookie is cleared
	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			found = true
			if cookie.Value != "" {
				t.Error("session_id cookie should be cleared")
			}
			if cookie.MaxAge != -1 {
				t.Error("session_id cookie MaxAge should be -1 to delete it")
			}
		}
	}
	if !found {
		t.Error("session_id cookie should be set to clear it")
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Logged out successfully" {
		t.Errorf("expected 'Logged out successfully', got %q", resp.Message)
	}
}

func TestLogout_NoSession(t *testing.T) {
	mockAuthService := &MockAuthService{
		DeleteSessionFunc: func(sessionID string) (int, error) {
			return http.StatusOK, nil
		},
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// Create a request without session cookie
	req := httptest.NewRequest("POST", "/logout", nil)

	rr := httptest.NewRecorder()
	authHandler.LogoutHandler(rr, req)

	// Should still return success (logout is idempotent)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Logged out successfully" {
		t.Errorf("expected 'Logged out successfully', got %q", resp.Message)
	}
}

// Test session invalidation after logout
func TestLogout_SessionInvalidation(t *testing.T) {
	sessionDeleted := false
	mockAuthService := &MockAuthService{
		DeleteSessionFunc: func(sessionID string) (int, error) {
			if sessionID == "valid-session-123" {
				sessionDeleted = true
				return http.StatusOK, nil
			}
			return 0, errors.New("session not found")
		},
		GetUserIDBySessionFunc: func(sessionID string) (int, error) {
			if sessionID == "valid-session-123" && !sessionDeleted {
				return 1, nil
			}
			return 0, errors.New("invalid or expired session")
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	// First, verify session is valid
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "valid-session-123",
	})

	rr := httptest.NewRecorder()
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
	})
	authMiddleware := authHandler.AuthMiddleware(protectedHandler)
	authMiddleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("Session should be valid before logout")
	}

	// Now logout
	logoutReq := httptest.NewRequest("POST", "/logout", nil)
	logoutReq.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "valid-session-123",
	})

	logoutRr := httptest.NewRecorder()
	authHandler.LogoutHandler(logoutRr, logoutReq)

	if logoutRr.Code != http.StatusOK {
		t.Error("Logout should succeed")
	}

	// Try to access protected route again with same session
	req2 := httptest.NewRequest("GET", "/protected", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "valid-session-123",
	})

	rr2 := httptest.NewRecorder()
	authMiddleware.ServeHTTP(rr2, req2)

	// Should now be unauthorized
	if rr2.Code != http.StatusUnauthorized {
		t.Error("Session should be invalid after logout")
	}
}
