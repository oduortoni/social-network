package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
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
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Unable to parse form: " + err.Error()})
		return
	}

	var post models.Post
	post.Content = r.FormValue("content") // Assuming post content is sent as a form value
	post.Privacy = r.FormValue("privacy")

	// Get user ID from context
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}
	post.UserID = userID

	// Handle optional image upload using the helper
	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		utils.RespondJSON(w, status, utils.Response{Message: err.Error()})
		return
	}

	// Handle private post viewers
	var viewerIDs []int64
	if post.Privacy == "private" {
		viewersParam := r.FormValue("viewers")
		if viewersParam != "" {
			// Parse comma-separated viewer IDs
			viewerIDStrings := strings.Split(viewersParam, ",")
			for _, idStr := range viewerIDStrings {
				idStr = strings.TrimSpace(idStr)
				if idStr != "" {
					id, err := strconv.ParseInt(idStr, 10, 64)
					if err != nil {
						utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid viewer ID: " + idStr})
						return
					}
					viewerIDs = append(viewerIDs, id)
				}
			}
		}
	}

	// Create post with viewers if it's private
	var id int64
	if post.Privacy == "private" && len(viewerIDs) > 0 {
		id, err = h.PostService.CreatePostWithViewers(&post, imageData, imageMimeType, viewerIDs)
	} else {
		id, err = h.PostService.CreatePost(&post, imageData, imageMimeType)
	}

	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	post.ID = id

	utils.RespondJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(20 << 20) // 20 MB limit for multipart form
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Unable to parse form: " + err.Error()})
		return
	}

	var comment models.Comment
	comment.Content = r.FormValue("content") // Assuming post content is sent as a form value

	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}
	comment.PostID = postID

	// Get user ID from context
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}
	comment.UserID = userID

	// Handle optional image upload using the helper
	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		utils.RespondJSON(w, status, utils.Response{Message: err.Error()})
		return
	}

	id, err := h.PostService.CreateComment(&comment, imageData, imageMimeType)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	comment.ID = id
	comment.CreatedAt = time.Now() // Set the created_at timestamp

	utils.RespondJSON(w, http.StatusCreated, comment)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	post, err := h.PostService.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondJSON(w, http.StatusNotFound, utils.Response{Message: "Post not found"})
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Internal server error"})
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	posts, err := h.PostService.GetPosts(userID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Internal server error"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	comments, err := h.PostService.GetCommentsByPostID(postID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Internal server error"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, comments)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	// Parse multipart form data
	err = r.ParseMultipartForm(20 << 20) // 20 MB limit for multipart form
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Unable to parse form: " + err.Error()})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	// Get the updated content
	content := r.FormValue("content")
	if content == "" {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Content is required"})
		return
	}

	// Handle optional image upload using the helper
	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		utils.RespondJSON(w, status, utils.Response{Message: err.Error()})
		return
	}

	// Update the post
	updatedPost, err := h.PostService.UpdatePost(postID, userID, content, imageData, imageMimeType)
	if err != nil {
		if err.Error() == "unauthorized" {
			utils.RespondJSON(w, http.StatusForbidden, utils.Response{Message: "You can only edit your own posts"})
		} else if err.Error() == "post not found" {
			utils.RespondJSON(w, http.StatusNotFound, utils.Response{Message: "Post not found"})
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, updatedPost)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	err = h.PostService.DeletePost(postID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		} else if err == sql.ErrNoRows {
			utils.RespondJSON(w, http.StatusNotFound, utils.Response{Message: "Post not found"})
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Internal server error"})
		}
		return
	}

	utils.RespondJSON(w, http.StatusNoContent, utils.Response{Message: "Post deleted successfully"})
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

func (h *PostHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	// Get search query from URL parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Search query is required"})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	users, err := h.PostService.SearchUsers(query, userID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to search users"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}
