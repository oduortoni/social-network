package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockPostService is a mock implementation of the PostService for testing.
type MockPostService struct {
	CreatePostFunc func(post *models.Post, imageData []byte, imageMimeType string) (int64, error)
}

func (s *MockPostService) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	return s.CreatePostFunc(post, imageData, imageMimeType)
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
		_, err = io.Copy(fw, bytes.NewReader(make([]byte, 10*1024*1024))) // 10MB dummy image
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

		ctx := context.WithValue(req.Context(), "userID", int64(1))
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

		ctx := context.WithValue(req.Context(), "userID", int64(1))
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

		ctx := context.WithValue(req.Context(), "userID", int64(1))
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
