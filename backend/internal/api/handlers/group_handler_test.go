package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockGroupService is a mock implementation of the GroupService for testing.
type MockGroupService struct {
	CreateGroupFunc func(group *models.Group) (*models.Group, error)
}

func (m *MockGroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	if m.CreateGroupFunc != nil {
		return m.CreateGroupFunc(group)
	}
	return group, nil
}

func TestCreateGroup(t *testing.T) {
	// Test case 1: Successful group creation
	t.Run("Successful group creation", func(t *testing.T) {
		mockService := &MockGroupService{
			CreateGroupFunc: func(group *models.Group) (*models.Group, error) {
				group.ID = 1
				return group, nil
			},
		}

		h := NewGroupHandler(mockService)

		group := &models.Group{
			Title:       "Test Group",
			Description: "This is a test group.",
			CreatorID:   1,
		}

		body, _ := json.Marshal(group)
		req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Add user ID to context
		ctx := context.WithValue(req.Context(), "userID", 1)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var createdGroup models.Group
		if err := json.NewDecoder(rr.Body).Decode(&createdGroup); err != nil {
			t.Fatal(err)
		}

		if createdGroup.ID != 1 {
			t.Errorf("handler returned unexpected group ID: got %v want %v",
				createdGroup.ID, 1)
		}
	})
}
