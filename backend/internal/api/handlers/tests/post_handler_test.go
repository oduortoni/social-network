package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockPostService is a mock implementation of the PostService for testing.
type MockPostService struct {
	CreatePostFunc func(post *models.Post) (int64, error)
}

func (s *MockPostService) CreatePost(post *models.Post) (int64, error) {
	return s.CreatePostFunc(post)
}

func TestCreatePost(t *testing.T) {
	// Create a new mock post service
	mockPostService := &MockPostService{
		CreatePostFunc: func(post *models.Post) (int64, error) {
			return 1, nil
		},
	}

	// Create a new post handler with the mock service
	postHandler := handlers.NewPostHandler(mockPostService)

	// Create a new request
	post := models.Post{Content: "Test post"}
	body, _ := json.Marshal(post)
	req, err := http.NewRequest("POST", "/posts", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}

	// Add user ID to context
	ctx := context.WithValue(req.Context(), "userID", int64(1))
	req = req.WithContext(ctx)

	// Create a new recorder
	rr := httptest.NewRecorder()

	// Call the handler
	postHandler.CreatePost(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body
	var createdPost models.Post
	if err := json.NewDecoder(rr.Body).Decode(&createdPost); err != nil {
		t.Fatal(err)
	}

	if createdPost.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			createdPost.ID, 1)
	}
}
