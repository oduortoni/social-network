package service

import (
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockPostStorePagination for testing pagination
type MockPostStorePagination struct {
	posts []models.Post
}

func (m *MockPostStorePagination) CreatePost(post *models.Post) (int64, error) { return 0, nil }
func (m *MockPostStorePagination) CreateComment(comment *models.Comment) (int64, error) { return 0, nil }
func (m *MockPostStorePagination) GetPostByID(id int64) (*models.Post, error) { return nil, nil }
func (m *MockPostStorePagination) GetPosts(userID int64) ([]*models.Post, error) { return nil, nil }
func (m *MockPostStorePagination) UpdatePost(postID int64, content, imagePath string) (*models.Post, error) { return nil, nil }
func (m *MockPostStorePagination) GetCommentsByPostID(postID, userID int64) ([]*models.Comment, error) { return nil, nil }
func (m *MockPostStorePagination) DeletePost(postID int64) error { return nil }
func (m *MockPostStorePagination) AddPostViewers(postID int64, viewerIDs []int64) error { return nil }
func (m *MockPostStorePagination) SearchUsers(query string, currentUserID int64) ([]*models.User, error) { return nil, nil }
func (m *MockPostStorePagination) UpdateComment(commentID int64, content, imagePath string) (*models.Comment, error) { return nil, nil }
func (m *MockPostStorePagination) DeleteComment(commentID int64) error { return nil }
func (m *MockPostStorePagination) GetCommentByID(commentID int64) (*models.Comment, error) { return nil, nil }

func (m *MockPostStorePagination) GetPostsPaginated(userID int64, limit, offset int) ([]*models.Post, error) {
	if limit == 0 {
		// Return all posts
		result := make([]*models.Post, len(m.posts))
		for i := range m.posts {
			result[i] = &m.posts[i]
		}
		return result, nil
	}

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

func (m *MockPostStorePagination) GetPostsCount(userID int64) (int, error) {
	return len(m.posts), nil
}

func TestGetPostsPaginated(t *testing.T) {
	// Create mock posts
	mockPosts := []models.Post{
		{ID: 1, Content: "Post 1"},
		{ID: 2, Content: "Post 2"},
		{ID: 3, Content: "Post 3"},
		{ID: 4, Content: "Post 4"},
		{ID: 5, Content: "Post 5"},
	}

	mockStore := &MockPostStorePagination{posts: mockPosts}
	service := NewPostService(mockStore)

	tests := []struct {
		name           string
		limit          int
		offset         int
		expectedCount  int
		expectedFirstID int64
	}{
		{"First page", 2, 0, 2, 1},
		{"Second page", 2, 2, 2, 3},
		{"Last page", 2, 4, 1, 5},
		{"Beyond range", 2, 10, 0, 0},
		{"All posts", 0, 0, 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := service.GetPostsPaginated(1, tt.limit, tt.offset)
			if err != nil {
				t.Errorf("GetPostsPaginated() error = %v", err)
				return
			}

			if len(posts) != tt.expectedCount {
				t.Errorf("GetPostsPaginated() got %d posts, want %d", len(posts), tt.expectedCount)
			}

			if tt.expectedCount > 0 && posts[0].ID != tt.expectedFirstID {
				t.Errorf("GetPostsPaginated() first post ID = %d, want %d", posts[0].ID, tt.expectedFirstID)
			}
		})
	}
}

func TestGetPostsCount(t *testing.T) {
	mockPosts := []models.Post{
		{ID: 1, Content: "Post 1"},
		{ID: 2, Content: "Post 2"},
		{ID: 3, Content: "Post 3"},
	}

	mockStore := &MockPostStorePagination{posts: mockPosts}
	service := NewPostService(mockStore)

	count, err := service.GetPostsCount(1)
	if err != nil {
		t.Errorf("GetPostsCount() error = %v", err)
		return
	}

	if count != 3 {
		t.Errorf("GetPostsCount() = %d, want 3", count)
	}
}