package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// MockPostServiceForPagination for testing pagination in handlers
type MockPostServiceForPagination struct {
	posts []models.Post
}

func (m *MockPostServiceForPagination) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	return 0, nil
}
func (m *MockPostServiceForPagination) CreatePostWithViewers(post *models.Post, imageData []byte, imageMimeType string, viewerIDs []int64) (int64, error) {
	return 0, nil
}
func (m *MockPostServiceForPagination) GetPostByID(id int64) (*models.Post, error) { return nil, nil }
func (m *MockPostServiceForPagination) UpdatePost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Post, error) {
	return nil, nil
}
func (m *MockPostServiceForPagination) CreateComment(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error) {
	return 0, nil
}
func (m *MockPostServiceForPagination) GetCommentsByPostID(postID, userID int64) ([]*models.Comment, error) {
	return nil, nil
}
func (m *MockPostServiceForPagination) DeletePost(postID, userID int64) error { return nil }
func (m *MockPostServiceForPagination) SearchUsers(query string, currentUserID int64) ([]*models.User, error) {
	return nil, nil
}
func (m *MockPostServiceForPagination) UpdateComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error) {
	return nil, nil
}
func (m *MockPostServiceForPagination) DeleteComment(commentID, userID int64) error { return nil }
func (m *MockPostServiceForPagination) GetCommentByID(commentID int64) (*models.Comment, error) {
	return nil, nil
}

func (m *MockPostServiceForPagination) GetPosts(userID int64) ([]*models.Post, error) {
	result := make([]*models.Post, len(m.posts))
	for i := range m.posts {
		result[i] = &m.posts[i]
	}
	return result, nil
}

func (m *MockPostServiceForPagination) GetPostsPaginated(userID int64, limit, offset int) ([]*models.Post, error) {
	start := offset
	end := offset + limit
	if start >= len(m.posts) {
		return []*models.Post{}, nil
	}
	if end > len(m.posts) {
		end = len(m.posts)
	}

	result := make([]*models.Post, end-start)
	for i := start; i < end; i++ {
		result[i-start] = &m.posts[i]
	}
	return result, nil
}

func (m *MockPostServiceForPagination) GetPostsCount(userID int64) (int, error) {
	return len(m.posts), nil
}

func TestGetPostsWithPagination(t *testing.T) {
	// Create mock posts
	mockPosts := []models.Post{
		{ID: 1, Content: "Post 1"},
		{ID: 2, Content: "Post 2"},
		{ID: 3, Content: "Post 3"},
		{ID: 4, Content: "Post 4"},
		{ID: 5, Content: "Post 5"},
	}

	mockService := &MockPostServiceForPagination{posts: mockPosts}
	handler := handlers.NewPostHandler(mockService)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkPagination bool
	}{
		{"No pagination params", "", http.StatusOK, false},
		{"First page", "?page=1&limit=2", http.StatusOK, true},
		{"Second page", "?page=2&limit=2", http.StatusOK, true},
		{"Invalid page", "?page=invalid&limit=2", http.StatusOK, true},
		{"Large limit", "?page=1&limit=100", http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/posts"+tt.queryParams, nil)
			
			// Add mock user context
			ctx := req.Context()
			ctx = utils.SetUserContext(ctx, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.GetPosts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetPosts() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkPagination {
				var response utils.PostsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
					return
				}

				if response.Pagination.CurrentPage == 0 {
					t.Error("Expected pagination metadata in response")
				}
			}
		})
	}
}