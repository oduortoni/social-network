package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/utils"
)

// MockUnfollowService is a mock implementation of the UnfollowServiceInterface for testing.
type MockUnfollowService struct {
	GetFollowConnectionIDFunc  func(followerID, followeeID int64) (int64, error)
	DeleteFollowConnectionFunc func(followConnectionID int64) error
}

func (s *MockUnfollowService) GetFollowConnectionID(followerID, followeeID int64) (int64, error) {
	if s.GetFollowConnectionIDFunc != nil {
		return s.GetFollowConnectionIDFunc(followerID, followeeID)
	}
	return 0, nil
}

func (s *MockUnfollowService) DeleteFollowConnection(followConnectionID int64) error {
	if s.DeleteFollowConnectionFunc != nil {
		return s.DeleteFollowConnectionFunc(followConnectionID)
	}
	return nil
}

func TestUnfollow_Success(t *testing.T) {
	mockUnfollowService := &MockUnfollowService{
		GetFollowConnectionIDFunc: func(followerID, followeeID int64) (int64, error) {
			if followerID == 1 && followeeID == 2 {
				return 123, nil // Mock connection ID
			}
			return 0, sql.ErrNoRows
		},
		DeleteFollowConnectionFunc: func(followConnectionID int64) error {
			if followConnectionID == 123 {
				return nil
			}
			return sql.ErrNoRows
		},
	}
	unfollowHandler := handlers.NewUnfollowHandler(mockUnfollowService)

	followee := models.Followee{FolloweeId: 2}
	body, _ := json.Marshal(followee)
	req := httptest.NewRequest("DELETE", "/unfollow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context (simulating authenticated user)
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	unfollowHandler.Unfollow(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Successfully unfollowed user"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestUnfollow_FollowRelationshipNotFound(t *testing.T) {
	mockUnfollowService := &MockUnfollowService{
		GetFollowConnectionIDFunc: func(followerID, followeeID int64) (int64, error) {
			return 0, sql.ErrNoRows // No follow relationship exists
		},
	}
	unfollowHandler := handlers.NewUnfollowHandler(mockUnfollowService)

	followee := models.Followee{FolloweeId: 2}
	body, _ := json.Marshal(followee)
	req := httptest.NewRequest("DELETE", "/unfollow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	unfollowHandler.Unfollow(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Follow relationship not found"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestUnfollow_InvalidJSON(t *testing.T) {
	mockUnfollowService := &MockUnfollowService{}
	unfollowHandler := handlers.NewUnfollowHandler(mockUnfollowService)

	// Invalid JSON
	req := httptest.NewRequest("DELETE", "/unfollow", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	unfollowHandler.Unfollow(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Invalid JSON format"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestUnfollow_Unauthorized(t *testing.T) {
	mockUnfollowService := &MockUnfollowService{}
	unfollowHandler := handlers.NewUnfollowHandler(mockUnfollowService)

	followee := models.Followee{FolloweeId: 2}
	body, _ := json.Marshal(followee)
	req := httptest.NewRequest("DELETE", "/unfollow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Don't add user ID to context (simulating unauthenticated user)

	rr := httptest.NewRecorder()
	unfollowHandler.Unfollow(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "User not found in context"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}
