package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/internal/utils"
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
	userID, ok := r.Context().Value(User_id).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	post.UserID = userID

	// Handle image upload
	file, handler, err := r.FormFile("image")
	var imageData []byte
	var imageMimeType string

	if err == nil { // No error means an image was provided
		defer file.Close()

		if handler.Size > 20*1024*1024 { // 20 MB limit
			http.Error(w, "Image size exceeds 20MB limit", http.StatusBadRequest)
			return
		}

		// Create a TeeReader to read from the file and also pass to checkImageSignature
		var buf bytes.Buffer
		teeReader := io.TeeReader(file, &buf)

		// Perform image signature check
		if _, err := utils.DetectImageFormat(teeReader); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Read the entire image data from the buffer and the remaining file content
		imageData, err = io.ReadAll(io.MultiReader(&buf, file))
		if err != nil {
			http.Error(w, "Failed to read image data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		imageMimeType = handler.Header.Get("Content-Type")
	} else if err != http.ErrMissingFile {
		// Other errors during FormFile processing
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusInternalServerError)
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
	json.NewEncoder(w).Encode(post)
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
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(User_id).(int64)
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
	json.NewEncoder(w).Encode(posts)
}
