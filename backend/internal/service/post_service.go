package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/internal/utils"
)



type PostService struct {
	PostStore *store.PostStore
}

func NewPostService(ps *store.PostStore) *PostService {
	return &PostService{PostStore: ps}
}

func (s *PostService) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	if len(imageData) > 0 {
		// Perform image signature check and get detected format
		detectedFormat, err := utils.DetectImageFormat(bytes.NewReader(imageData))
		if err != nil {
			return 0, fmt.Errorf("image signature check failed: %w", err)
		}

		// Determine file extension from detected format (more reliable than MIME type)
		var extension string
		switch detectedFormat {
		case utils.JPEG:
			extension = ".jpg"
		case utils.PNG:
			extension = ".png"
		case utils.GIF:
			extension = ".gif"
		case utils.WebP:
			extension = ".webp"
		case utils.BMP:
			extension = ".bmp"
		case utils.TIFF:
			extension = ".tiff"
		default:
			return 0, fmt.Errorf("unsupported image format: %s", detectedFormat)
		}

		// Generate a unique filename
		uuid := uuid.New()
		imageFileName := fmt.Sprintf("%s%s", uuid.String(), extension)
		imagePath := filepath.Join("PostImages", imageFileName)

		// Create the directory if it doesn't exist
		err = os.MkdirAll("PostImages", os.ModePerm)
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
