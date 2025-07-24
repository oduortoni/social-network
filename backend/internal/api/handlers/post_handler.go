package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/utils"
)

type PostHandler struct {
	PostService service.PostServiceInterface
}

func NewPostHandler(ps service.PostServiceInterface) *PostHandler {
	return &PostHandler{PostService: ps}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(20 << 20) // 20 MB limit for multipart form
	if err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var post models.Post
	post.Content = r.FormValue("content") // Assuming post content is sent as a form value
	post.Privacy = r.FormValue("privacy")

	// Get user ID from context
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	post.UserID = userID

	// Handle optional image upload using the helper
	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	id, err := h.PostService.CreatePost(&post, imageData, imageMimeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(20 << 20) // 20 MB limit for multipart form
	if err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var comment models.Comment
	comment.Content = r.FormValue("content") // Assuming post content is sent as a form value

	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	comment.PostID = postID

	// Get user ID from context
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	comment.UserID = userID

	// Handle optional image upload using the helper
	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	id, err := h.PostService.CreateComment(&comment, imageData, imageMimeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comment.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := h.PostService.GetFeed(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(posts)
}

// handleImageUpload processes an optional image from a multipart form.
// It returns the image data, its MIME type, an appropriate HTTP status code for errors, and any error encountered.
func handleImageUpload(r *http.Request) (imageData []byte, imageMimeType string, status int, err error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, "", 0, nil // No image provided, not an error
		}
		// Other errors during FormFile processing
		return nil, "", http.StatusInternalServerError, fmt.Errorf("error retrieving the file: %w", err)
	}
	defer file.Close()

	// 20 MB limit
	const maxImageSize = 20 << 20
	if handler.Size > maxImageSize {
		return nil, "", http.StatusBadRequest, fmt.Errorf("image size exceeds 20MB limit")
	}

	// Use a TeeReader to read from the file for both signature check and full read
	var buf bytes.Buffer
	teeReader := io.TeeReader(file, &buf)

	// Perform image signature check
	if _, err := utils.DetectImageFormat(teeReader); err != nil {
		return nil, "", http.StatusBadRequest, err
	}

	// Read the entire image data from the buffer and the remaining file content
	imageData, err = io.ReadAll(io.MultiReader(&buf, file))
	if err != nil {
		return nil, "", http.StatusInternalServerError, fmt.Errorf("failed to read image data: %w", err)
	}
	imageMimeType = handler.Header.Get("Content-Type")

	return imageData, imageMimeType, 0, nil
}
