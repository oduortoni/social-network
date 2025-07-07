package service

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type PostService struct {
	PostStore *store.PostStore
}

func NewPostService(ps *store.PostStore) *PostService {
	return &PostService{PostStore: ps}
}

func (s *PostService) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	if len(imageData) > 0 {
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
