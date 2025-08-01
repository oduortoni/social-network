package service

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockPostStore is a mock implementation of the PostStore for testing.
type MockPostStore struct {
	CreatePostFunc          func(post *models.Post) (int64, error)
	GetPostsFunc            func(userID int64) ([]*models.Post, error)
	GetCommentsByPostIDFunc func(postID, userID int64) ([]*models.Comment, error)
	GetPostByIDFunc         func(id int64) (*models.Post, error)
	DeletePostFunc          func(postID int64) error
	AddPostViewersFunc      func(postID int64, viewerIDs []int64) error
}

func (s *MockPostStore) CreatePost(post *models.Post) (int64, error) {
	if s.CreatePostFunc != nil {
		return s.CreatePostFunc(post)
	}
	return 1, nil
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
	if s.AddPostViewersFunc != nil {
		return s.AddPostViewersFunc(postID, viewerIDs)
	}
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

func TestCreatePost(t *testing.T) {
	// Test case 1: Successful post creation
	t.Run("Successful post creation", func(t *testing.T) {
		mockStore := &MockPostStore{
			CreatePostFunc: func(post *models.Post) (int64, error) {
				if post.Content == "" {
					t.Errorf("expected post content to be non-empty")
				}
				return 1, nil
			},
		}
		postService := NewPostService(mockStore)

		post := &models.Post{
			Content: "Test Post",
		}
		_, err := postService.CreatePost(post, nil, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	// Test case 2: Post creation with no content
	t.Run("Post creation with no content", func(t *testing.T) {
		mockStore := &MockPostStore{}
		postService := NewPostService(mockStore)

		post := &models.Post{
			Content: "",
		}
		_, err := postService.CreatePost(post, nil, "")
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
		if err.Error() != "post content is required" {
			t.Fatalf("expected error 'post content is required', got %v", err)
		}
	})

	// Test case 3: Post creation with image
	t.Run("Post creation with image", func(t *testing.T) {
		mockStore := &MockPostStore{
			CreatePostFunc: func(post *models.Post) (int64, error) {
				if post.Image == "" {
					t.Errorf("expected post image to be non-empty")
				}
				return 1, nil
			},
		}
		postService := NewPostService(mockStore)

		post := &models.Post{
			Content: "Test Post",
		}
		// a simple fake image data
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		imagePath := filepath.Join(wd, "..", "..", "..", "entity_relation_diagram.png")
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			t.Fatalf("failed to read image file: %v", err)
		}
		_, err = postService.CreatePost(post, imageData, "image/png")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	// Test case 4: Post creation with image and no content
	t.Run("Post creation with image and no content", func(t *testing.T) {
		mockStore := &MockPostStore{}
		postService := NewPostService(mockStore)

		post := &models.Post{
			Content: "",
		}
		// a simple fake image data
		imageData := []byte("fake-image-data")
		_, err := postService.CreatePost(post, imageData, "image/jpeg")
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
		if err.Error() != "post content is required" {
			t.Fatalf("expected error 'post content is required', got %v", err)
		}
	})

	// Test case 5: Private post creation with viewers
	t.Run("Private post creation with viewers", func(t *testing.T) {
		var capturedPostID int64
		var capturedViewerIDs []int64

		mockStore := &MockPostStore{
			CreatePostFunc: func(post *models.Post) (int64, error) {
				return 1, nil
			},
			AddPostViewersFunc: func(postID int64, viewerIDs []int64) error {
				capturedPostID = postID
				capturedViewerIDs = viewerIDs
				return nil
			},
		}
		postService := NewPostService(mockStore)

		post := &models.Post{
			Content: "Private Post",
			Privacy: "private",
		}
		viewerIDs := []int64{101, 102}
		postID, err := postService.CreatePostWithViewers(post, nil, "", viewerIDs)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if postID != 1 {
			t.Fatalf("expected postID to be 1, got %d", postID)
		}
		if capturedPostID != 1 {
			t.Fatalf("expected captured postID to be 1, got %d", capturedPostID)
		}
		if len(capturedViewerIDs) != 2 {
			t.Fatalf("expected 2 viewerIDs, got %d", len(capturedViewerIDs))
		}
	})
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
