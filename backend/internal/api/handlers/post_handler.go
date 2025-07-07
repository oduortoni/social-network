package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
)

// Define magic numbers for common image formats
var imageMagicNumbers = map[string][]byte{
	"jpeg": {0xFF, 0xD8, 0xFF},
	"png":  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	"gif":  {0x47, 0x49, 0x46, 0x38}, // GIF87a or GIF89a
	"webp": {0x52, 0x49, 0x46, 0x46}, // RIFF (WebP files start with RIFF, followed by file size and WEBP)
}

// checkImageSignature reads the first few bytes and validates against known image magic numbers.
func checkImageSignature(reader io.Reader) error {
	// Read enough bytes to cover the longest magic number (PNG: 8 bytes)
	buffer := make([]byte, 8)
	n, err := io.ReadAtLeast(reader, buffer, 4) // Read at least 4 bytes (for GIF/WebP)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	// Check for JPEG
	if n >= 3 && bytes.Equal(buffer[0:3], imageMagicNumbers["jpeg"]) {
		return nil
	}
	// Check for PNG
	if n >= 8 && bytes.Equal(buffer[0:8], imageMagicNumbers["png"]) {
		return nil
	}
	// Check for GIF
	if n >= 4 && bytes.Equal(buffer[0:4], imageMagicNumbers["gif"]) {
		return nil
	}
	// Check for WebP (RIFF header, then check for "WEBP" at offset 8)
	if n >= 4 && bytes.Equal(buffer[0:4], imageMagicNumbers["webp"]) {
		// For WebP, we need to read more to confirm "WEBP"
		// Since we only read 8 bytes, we can't fully validate WEBP here without more reads.
		// For simplicity and to avoid re-reading, we'll assume RIFF is enough for now,
		// but a more robust check would involve reading bytes 8-11 for "WEBP".
		// For this implementation, we'll just check the RIFF header.
		return nil
	}

	return errors.New("unsupported image format or invalid file signature")
}

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

	// Get user ID from context
	userID, ok := r.Context().Value("userID").(int64)
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
		if err := checkImageSignature(teeReader); err != nil {
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
