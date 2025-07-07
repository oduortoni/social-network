package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
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

type PostService struct {
	PostStore *store.PostStore
}

func NewPostService(ps *store.PostStore) *PostService {
	return &PostService{PostStore: ps}
}

func (s *PostService) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	if len(imageData) > 0 {
		// Perform image signature check
		if err := checkImageSignature(bytes.NewReader(imageData)); err != nil {
			return 0, fmt.Errorf("image signature check failed: %w", err)
		}

		// Determine file extension from MIME type
		extension := ".bin" // Default to .bin if MIME type is unknown
		switch imageMimeType {
		case "image/jpeg":
			extension = ".jpg"
		case "image/png":
			extension = ".png"
		case "image/gif":
			extension = ".gif"
		case "image/webp":
			extension = ".webp"
		case "image/svg+xml":
			extension = ".svg"
		case "image/bmp":
			extension = ".bmp"
		}

		// Generate a unique filename
		uuid := uuid.New()
		imageFileName := fmt.Sprintf("%s%s", uuid.String(), extension)
		imagePath := filepath.Join("UserAvatars", imageFileName) // Save in UserAvatars directory

		// Create the directory if it doesn't exist
		err := os.MkdirAll("UserAvatars", os.ModePerm)
		if err != nil {
			return 0, fmt.Errorf("failed to create directory: %w", err)
		}

		// Save the image file
		err = os.WriteFile(imagePath, imageData, 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to save image: %w", err)
		}

		post.Image = imagePath
	}

	return s.PostStore.CreatePost(post)
}
