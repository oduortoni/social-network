package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// MockPostService is a mock implementation of the PostService for testing.
type MockPostService struct {
	CreatePostFunc          func(post *models.Post, imageData []byte, imageMimeType string) (int64, error)
	GetPostByIDFunc         func(id int64) (*models.Post, error)
	GetPostsFunc            func(userID int64) ([]*models.Post, error)
	CreateCommentFunc       func(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error)
	GetCommentsByPostIDFunc func(postID, userID int64) ([]*models.Comment, error)
	DeletePostFunc          func(postID, userID int64) error
	UpdateCommentFunc       func(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error)
	DeleteCommentFunc       func(commentID, userID int64) error
	GetCommentByIDFunc      func(commentID int64) (*models.Comment, error)
}

func (s *MockPostService) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	return s.CreatePostFunc(post, imageData, imageMimeType)
}

func (s *MockPostService) CreateComment(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error) {
	if s.CreateCommentFunc != nil {
		return s.CreateCommentFunc(comment, imageData, imageMimeType)
	}
	return 0, fmt.Errorf("CreateCommentFunc not implemented")
}

func (s *MockPostService) GetPostByID(id int64) (*models.Post, error) {
	return s.GetPostByIDFunc(id)
}

func (s *MockPostService) GetPosts(userID int64) ([]*models.Post, error) {
	if s.GetPostsFunc != nil {
		return s.GetPostsFunc(userID)
	}
	return nil, fmt.Errorf("GetPostsFunc not implemented")
}

func (s *MockPostService) GetCommentsByPostID(postID, userID int64) ([]*models.Comment, error) {
	if s.GetCommentsByPostIDFunc != nil {
		return s.GetCommentsByPostIDFunc(postID, userID)
	}
	return nil, fmt.Errorf("GetCommentsByPostIDFunc not implemented")
}

func (s *MockPostService) DeletePost(postID, userID int64) error {
	if s.DeletePostFunc != nil {
		return s.DeletePostFunc(postID, userID)
	}
	return fmt.Errorf("DeletePostFunc not implemented")
}

func (s *MockPostService) CreatePostWithViewers(post *models.Post, imageData []byte, imageMimeType string, viewerIDs []int64) (int64, error) {
	if s.CreatePostFunc != nil {
		return s.CreatePostFunc(post, imageData, imageMimeType)
	}
	return 0, fmt.Errorf("CreatePostFunc not implemented")
}

func (s *MockPostService) UpdatePost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Post, error) {
	return nil, nil
}

func (s *MockPostService) SearchUsers(query string, currentUserID int64) ([]*models.User, error) {
	return nil, nil
}

func (s *MockPostService) UpdateComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error) {
	if s.UpdateCommentFunc != nil {
		return s.UpdateCommentFunc(commentID, userID, content, imageData, imageMimeType)
	}
	return nil, fmt.Errorf("UpdateCommentFunc not implemented")
}

func (s *MockPostService) DeleteComment(commentID, userID int64) error {
	if s.DeleteCommentFunc != nil {
		return s.DeleteCommentFunc(commentID, userID)
	}
	return fmt.Errorf("DeleteCommentFunc not implemented")
}

func (s *MockPostService) GetCommentByID(commentID int64) (*models.Comment, error) {
	if s.GetCommentByIDFunc != nil {
		return s.GetCommentByIDFunc(commentID)
	}
	return nil, fmt.Errorf("GetCommentByIDFunc not implemented")
}

func TestDeletePost(t *testing.T) {
	// Test case 1: Successful deletion
	t.Run("Successful deletion", func(t *testing.T) {
		mockPostService := &MockPostService{
			DeletePostFunc: func(postID, userID int64) error {
				if postID != 1 || userID != 100 {
					t.Errorf("unexpected input to DeletePost: postID=%d, userID=%d", postID, userID)
				}
				return nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/posts/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeletePost(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})

	// Test case 2: Invalid post ID
	t.Run("Invalid post ID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("DELETE", "/posts/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "invalid")

		rr := httptest.NewRecorder()
		postHandler.DeletePost(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 3: Unauthorized (missing userID in context)
	t.Run("Unauthorized missing userID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("DELETE", "/posts/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		rr := httptest.NewRecorder()
		postHandler.DeletePost(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Unauthorized (service returns unauthorized error)
	t.Run("Unauthorized service error", func(t *testing.T) {
		mockPostService := &MockPostService{
			DeletePostFunc: func(postID, userID int64) error {
				return fmt.Errorf("unauthorized")
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/posts/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeletePost(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	// Test case 5: Post not found (service returns sql.ErrNoRows)
	t.Run("Post not found service error", func(t *testing.T) {
		mockPostService := &MockPostService{
			DeletePostFunc: func(postID, userID int64) error {
				return sql.ErrNoRows
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/posts/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeletePost(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 6: Internal server error (service returns generic error)
	t.Run("Internal server error", func(t *testing.T) {
		mockPostService := &MockPostService{
			DeletePostFunc: func(postID, userID int64) error {
				return fmt.Errorf("database error")
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/posts/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeletePost(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestCreatePost(t *testing.T) {
	// Test case 1: Successful post creation with image
	t.Run("Successful post creation with image", func(t *testing.T) {
		mockPostService := &MockPostService{
			CreatePostFunc: func(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
				if post.Content != "Test post with image" || imageData == nil || imageMimeType != "image/jpeg" {
					t.Errorf("unexpected input to CreatePost: content=%s, imageData present=%t, mimeType=%s", post.Content, imageData != nil, imageMimeType)
				}
				return 1, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		fw, err := w.CreatePart(map[string][]string{
			"Content-Disposition": {fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "test_image.jpeg")},
			"Content-Type":        {"image/jpeg"},
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(fw, bytes.NewReader([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46})) // 10MB dummy image
		if err != nil {
			t.Fatal(err)
		}
		_ = w.WriteField("content", "Test post with image")
		w.Close()

		req, err := http.NewRequest("POST", "/posts", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.CreatePost(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	// Test case 2: Image size exceeds limit
	t.Run("Image size exceeds limit", func(t *testing.T) {
		mockPostService := &MockPostService{
			CreatePostFunc: func(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
				return 0, nil // Should not be called if validation fails
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		fw, err := w.CreateFormFile("image", "large_image.png")
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(fw, bytes.NewReader(make([]byte, 21*1024*1024))) // 21MB dummy image
		if err != nil {
			t.Fatal(err)
		}
		_ = w.WriteField("content", "Test post with large image")
		w.Close()

		req, err := http.NewRequest("POST", "/posts", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.CreatePost(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		if !strings.Contains(rr.Body.String(), "image size exceeds 20MB limit") {
			t.Errorf("handler returned unexpected error message: %s", rr.Body.String())
		}
	})

	// Test case 3: No image provided
	t.Run("No image provided", func(t *testing.T) {
		mockPostService := &MockPostService{
			CreatePostFunc: func(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
				if post.Content != "Test post without image" || imageData != nil || imageMimeType != "" {
					t.Errorf("unexpected input to CreatePost: content=%s, imageData present=%t, mimeType=%s", post.Content, imageData != nil, imageMimeType)
				}
				return 1, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Test post without image")
		w.Close()

		req, err := http.NewRequest("POST", "/posts", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.CreatePost(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	// Test case 4: Unauthorized (missing userID in context)
	t.Run("Unauthorized missing userID", func(t *testing.T) {
		mockPostService := &MockPostService{
			CreatePostFunc: func(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
				return 0, nil // Should not be called
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Test post unauthorized")
		w.Close()

		req, err := http.NewRequest("POST", "/posts", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		rr := httptest.NewRecorder()
		postHandler.CreatePost(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
		if !strings.Contains(rr.Body.String(), "Unauthorized") {
			t.Errorf("handler returned unexpected error message: %s", rr.Body.String())
		}
	})
}

func TestGetPostByID(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Successful retrieval", func(t *testing.T) {
		mockPostService := &MockPostService{
			GetPostByIDFunc: func(id int64) (*models.Post, error) {
				if id != 1 {
					t.Errorf("unexpected post ID: got %v want %v", id, 1)
				}
				return &models.Post{ID: 1, Content: "Test Post"}, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("GET", "/posts/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		rr := httptest.NewRecorder()
		postHandler.GetPostByID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var post models.Post
		if err := json.NewDecoder(rr.Body).Decode(&post); err != nil {
			t.Fatal(err)
		}

		if post.ID != 1 || post.Content != "Test Post" {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})

	// Test case 2: Post not found
	t.Run("Post not found", func(t *testing.T) {
		mockPostService := &MockPostService{
			GetPostByIDFunc: func(id int64) (*models.Post, error) {
				return nil, sql.ErrNoRows
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("GET", "/posts/2", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "2")

		rr := httptest.NewRecorder()
		postHandler.GetPostByID(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 3: Invalid post ID
	t.Run("Invalid post ID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("GET", "/posts/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "invalid")

		rr := httptest.NewRecorder()
		postHandler.GetPostByID(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestGetPosts(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Successful retrieval", func(t *testing.T) {
		mockPostService := &MockPostService{
			GetPostsFunc: func(userID int64) ([]*models.Post, error) {
				if userID != 1 {
					t.Errorf("unexpected user ID: got %v want %v", userID, 1)
				}
				return []*models.Post{{ID: 1, Content: "Test Post"}}, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.GetPosts(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var posts []*models.Post
		if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
			t.Fatal(err)
		}

		if len(posts) != 1 || posts[0].ID != 1 || posts[0].Content != "Test Post" {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})

	// Test case 2: Unauthorized
	t.Run("Unauthorized", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		postHandler.GetPosts(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}

func TestGetCommentsByPostID(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Successful retrieval", func(t *testing.T) {
		mockPostService := &MockPostService{
			GetCommentsByPostIDFunc: func(postID, userID int64) ([]*models.Comment, error) {
				if postID != 1 {
					t.Errorf("unexpected post ID: got %v want %v", postID, 1)
				}
				return []*models.Comment{{ID: 1, Content: "Test Comment"}}, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("GET", "/posts/1/comments", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.GetCommentsByPostID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var comments []*models.Comment
		if err := json.NewDecoder(rr.Body).Decode(&comments); err != nil {
			t.Fatal(err)
		}

		if len(comments) != 1 || comments[0].ID != 1 || comments[0].Content != "Test Comment" {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})

	// Test case 2: Invalid post ID
	t.Run("Invalid post ID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("GET", "/posts/invalid/comments", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "invalid")

		rr := httptest.NewRecorder()
		postHandler.GetCommentsByPostID(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestCreateComment(t *testing.T) {
	// Test case 1: Successful comment creation
	t.Run("Successful comment creation", func(t *testing.T) {
		mockPostService := &MockPostService{
			CreateCommentFunc: func(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error) {
				if comment.Content != "Test comment" || comment.PostID != 1 {
					t.Errorf("unexpected input to CreateComment: content=%s, postID=%d", comment.Content, comment.PostID)
				}
				return 1, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Test comment")
		w.Close()

		req, err := http.NewRequest("POST", "/posts/1/comments", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.CreateComment(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	// Test case 2: Invalid post ID
	t.Run("Invalid post ID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("POST", "/posts/invalid/comments", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("postId", "invalid")

		rr := httptest.NewRecorder()
		postHandler.CreateComment(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 3: Unauthorized
	t.Run("Unauthorized", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Test comment")
		w.Close()

		req, err := http.NewRequest("POST", "/posts/1/comments", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("postId", "1")

		rr := httptest.NewRecorder()
		postHandler.CreateComment(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Service layer error
	t.Run("Service layer error", func(t *testing.T) {
		mockPostService := &MockPostService{
			CreateCommentFunc: func(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error) {
				return 0, fmt.Errorf("service error")
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Test comment")
		w.Close()

		req, err := http.NewRequest("POST", "/posts/1/comments", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("postId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.CreateComment(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
func TestUpdateComment(t *testing.T) {
	// Test case 1: Successful comment update
	t.Run("Successful comment update", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.UpdateCommentFunc = func(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error) {
			if commentID != 1 || userID != 100 || content != "Updated comment content" {
				t.Errorf("unexpected input to UpdateComment: commentID=%d, userID=%d, content=%s", commentID, userID, content)
			}
			return &models.Comment{ID: 1, Content: "Updated comment content", UserID: 100, IsEdited: true}, nil
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Updated comment content")
		w.Close()

		req, err := http.NewRequest("PUT", "/comments/1", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var comment models.Comment
		if err := json.NewDecoder(rr.Body).Decode(&comment); err != nil {
			t.Fatal(err)
		}

		if comment.ID != 1 || comment.Content != "Updated comment content" || !comment.IsEdited {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})

	// Test case 2: Invalid comment ID
	t.Run("Invalid comment ID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("PUT", "/comments/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "invalid")

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 3: Unauthorized (missing userID in context)
	t.Run("Unauthorized missing userID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Updated comment")
		w.Close()

		req, err := http.NewRequest("PUT", "/comments/1", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("commentId", "1")

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Missing content
	t.Run("Missing content", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		w.Close()

		req, err := http.NewRequest("PUT", "/comments/1", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		if !strings.Contains(rr.Body.String(), "Content is required") {
			t.Errorf("handler returned unexpected error message: %s", rr.Body.String())
		}
	})

	// Test case 5: Unauthorized (service returns unauthorized error)
	t.Run("Unauthorized service error", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.UpdateCommentFunc = func(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error) {
			return nil, fmt.Errorf("unauthorized")
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Updated comment")
		w.Close()

		req, err := http.NewRequest("PUT", "/comments/1", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
		}

		if !strings.Contains(rr.Body.String(), "You can only edit your own comments") {
			t.Errorf("handler returned unexpected error message: %s", rr.Body.String())
		}
	})

	// Test case 6: Comment not found
	t.Run("Comment not found", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.UpdateCommentFunc = func(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error) {
			return nil, fmt.Errorf("comment not found")
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Updated comment")
		w.Close()

		req, err := http.NewRequest("PUT", "/comments/1", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 7: Internal server error
	t.Run("Internal server error", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.UpdateCommentFunc = func(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error) {
			return nil, fmt.Errorf("database error")
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_ = w.WriteField("content", "Updated comment")
		w.Close()

		req, err := http.NewRequest("PUT", "/comments/1", &b)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.UpdateComment(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestDeleteComment(t *testing.T) {
	// Test case 1: Successful deletion
	t.Run("Successful deletion", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.DeleteCommentFunc = func(commentID, userID int64) error {
			if commentID != 1 || userID != 100 {
				t.Errorf("unexpected input to DeleteComment: commentID=%d, userID=%d", commentID, userID)
			}
			return nil
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/comments/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeleteComment(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})

	// Test case 2: Invalid comment ID
	t.Run("Invalid comment ID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("DELETE", "/comments/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "invalid")

		rr := httptest.NewRecorder()
		postHandler.DeleteComment(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 3: Unauthorized (missing userID in context)
	t.Run("Unauthorized missing userID", func(t *testing.T) {
		postHandler := handlers.NewPostHandler(nil) // No service needed for this test

		req, err := http.NewRequest("DELETE", "/comments/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "1")

		rr := httptest.NewRecorder()
		postHandler.DeleteComment(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Unauthorized (service returns unauthorized error)
	t.Run("Unauthorized service error", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.DeleteCommentFunc = func(commentID, userID int64) error {
			return fmt.Errorf("unauthorized")
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/comments/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeleteComment(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
		}

		if !strings.Contains(rr.Body.String(), "You can only delete your own comments") {
			t.Errorf("handler returned unexpected error message: %s", rr.Body.String())
		}
	})

	// Test case 5: Comment not found (service returns sql.ErrNoRows)
	t.Run("Comment not found service error", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.DeleteCommentFunc = func(commentID, userID int64) error {
			return sql.ErrNoRows
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/comments/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeleteComment(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 6: Internal server error (service returns generic error)
	t.Run("Internal server error", func(t *testing.T) {
		mockPostService := &MockPostService{}
		mockPostService.DeleteCommentFunc = func(commentID, userID int64) error {
			return fmt.Errorf("database error")
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("DELETE", "/comments/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("commentId", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.DeleteComment(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
