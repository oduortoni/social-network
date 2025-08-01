package service

import (
	"fmt"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockPostStore is a mock implementation of the PostStore for testing.
type MockPostStore struct {
	GetPostsFunc            func(userID int64) ([]*models.Post, error)
	GetCommentsByPostIDFunc func(postID, userID int64) ([]*models.Comment, error)
	GetPostByIDFunc         func(id int64) (*models.Post, error)
	DeletePostFunc          func(postID int64) error
}

func (s *MockPostStore) CreatePost(post *models.Post) (int64, error) {
	return 0, nil
}

func (s *MockPostStore) CreateComment(comment *models.Comment) (int64, error) {
	return 0, nil
}

func (s *MockPostStore) GetPostByID(id int64) (*models.Post, error) {
	return s.GetPostByIDFunc(id)
}

func (s *MockPostStore) GetPosts(userID int64) ([]*models.Post, error) {
	return s.GetPostsFunc(userID)
}

func (s *MockPostStore) GetCommentsByPostID(postID, userID int64) ([]*models.Comment, error) {
	return s.GetCommentsByPostIDFunc(postID, userID)
}

func (s *MockPostStore) DeletePost(postID int64) error {
	return s.DeletePostFunc(postID)
}

func (s *MockPostStore) AddPostViewers(postID int64, viewerIDs []int64) error {
	return nil
}

func (s *MockPostStore) UpdatePost(postID int64, content, imagePath string) (*models.Post, error) {
	return nil, nil
}

func (s *MockPostStore) SearchUsers(query string, currentUserID int64) ([]*models.User, error) {
	return nil, nil
}

func (s *MockPostStore) UpdateComment(commentID int64, content, imagePath string) (*models.Comment, error) {
	return nil, nil
}

func (s *MockPostStore) DeleteComment(commentID int64) error {
	return nil
}

func (s *MockPostStore) GetCommentByID(commentID int64) (*models.Comment, error) {
	return nil, nil
}

func TestDeletePost(t *testing.T) {
	// Test case 1: Successful deletion
	t.Run("Successful deletion", func(t *testing.T) {
		mockStore := &MockPostStore{
			GetPostByIDFunc: func(id int64) (*models.Post, error) {
				return &models.Post{ID: 1, UserID: 100, Content: "Test Post"}, nil
			},
			DeletePostFunc: func(postID int64) error {
				if postID != 1 {
					t.Errorf("unexpected post ID for deletion: got %v want %v", postID, 1)
				}
				return nil
			},
		}
		postService := NewPostService(mockStore)

		err := postService.DeletePost(1, 100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	// Test case 2: Unauthorized deletion
	t.Run("Unauthorized deletion", func(t *testing.T) {
		mockStore := &MockPostStore{
			GetPostByIDFunc: func(id int64) (*models.Post, error) {
				return &models.Post{ID: 1, UserID: 100, Content: "Test Post"}, nil
			},
			DeletePostFunc: func(postID int64) error {
				t.Errorf("DeletePostFunc should not be called for unauthorized deletion")
				return nil
			},
		}
		postService := NewPostService(mockStore)

		err := postService.DeletePost(1, 101) // Different user ID
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("expected unauthorized error, got %v", err)
		}
	})

	// Test case 3: Post not found
	t.Run("Post not found", func(t *testing.T) {
		mockStore := &MockPostStore{
			GetPostByIDFunc: func(id int64) (*models.Post, error) {
				return nil, fmt.Errorf("post not found")
			},
			DeletePostFunc: func(postID int64) error {
				t.Errorf("DeletePostFunc should not be called if post not found")
				return nil
			},
		}
		postService := NewPostService(mockStore)

		err := postService.DeletePost(1, 100)
		if err == nil || err.Error() != "post not found" {
			t.Fatalf("expected 'post not found' error, got %v", err)
		}
	})
}
