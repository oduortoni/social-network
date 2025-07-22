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
	"github.com/tajjjjr/social-network/backend/utils"
)

// MockPostService is a mock implementation of the PostService for testing.
type MockPostService struct {
	CreatePostFunc    func(post *models.Post, imageData []byte, imageMimeType string) (int64, error)
	GetPostByIDFunc   func(id int64) (*models.Post, error)
	GetFeedFunc       func(userID int64) ([]*models.Post, error)
	CreateCommentFunc func(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error)
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

func (s *MockPostService) GetFeed(userID int64) ([]*models.Post, error) {
	if s.GetFeedFunc != nil {
		return s.GetFeedFunc(userID)
	}
	return nil, fmt.Errorf("GetFeedFunc not implemented")
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
			"Content-Disposition": []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "test_image.jpeg")},
			"Content-Type":        []string{"image/jpeg"},
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

		if !strings.Contains(rr.Body.String(), "Image size exceeds 20MB limit") {
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

func TestGetFeed(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Successful retrieval", func(t *testing.T) {
		mockPostService := &MockPostService{
			GetFeedFunc: func(userID int64) ([]*models.Post, error) {
				if userID != 1 {
					t.Errorf("unexpected user ID: got %v want %v", userID, 1)
				}
				return []*models.Post{{ID: 1, Content: "Test Post"}}, nil
			},
		}
		postHandler := handlers.NewPostHandler(mockPostService)

		req, err := http.NewRequest("GET", "/feed", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		postHandler.GetFeed(rr, req)

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

		req, err := http.NewRequest("GET", "/feed", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		postHandler.GetFeed(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
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
