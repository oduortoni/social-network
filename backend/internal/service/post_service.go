package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type PostService struct {
	PostStore store.PostStoreInterface
}

func NewPostService(ps store.PostStoreInterface) *PostService {
	return &PostService{PostStore: ps}
}

func (s *PostService) CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error) {
	if post.Content == "" {
		return 0, fmt.Errorf("post content is required")
	}
	if len(imageData) > 0 {
		// Perform image signature check and get detected format
		imagePath, err := s.saveImage(imageData, "posts")
		if err != nil {
			return 0, err
		}
		post.Image = imagePath
	}

	return s.PostStore.CreatePost(post)
}

func (s *PostService) CreateComment(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error) {
	if comment.Content == "" {
		return 0, fmt.Errorf("comment content is required")
	}
	if len(imageData) > 0 {
		// Perform image signature check and get detected format
		imagePath, err := s.saveImage(imageData, "comments")
		if err != nil {
			return 0, err
		}
		comment.Image = imagePath
	}
	return s.PostStore.CreateComment(comment)
}

func (s *PostService) GetPostByID(id int64) (*models.Post, error) {
	return s.PostStore.GetPostByID(id)
}

func (s *PostService) GetPosts(userID int64) ([]*models.Post, error) {
	return s.PostStore.GetPosts(userID)
}

func (s *PostService) GetCommentsByPostID(postID int64) ([]*models.Comment, error) {
	return s.PostStore.GetCommentsByPostID(postID)
}

func (s *PostService) DeletePost(postID, userID int64) error {
	post, err := s.PostStore.GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.PostStore.DeletePost(postID)
}

// saveImage handles the logic for validating, naming, and saving an uploaded image.
// It takes the image data and a sub-directory (e.g., "posts", "comments") to save the image in.
// It returns the saved file path or an error.
func (s *PostService) saveImage(imageData []byte, subDir string) (string, error) {
	// Perform image signature check and get detected format
	detectedFormat, err := utils.DetectImageFormat(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("image signature check failed: %w", err)
	}
	extension, ok := formatToExtension(detectedFormat)
	if !ok {
		return "", fmt.Errorf("unsupported image format: %s", detectedFormat)
	}
	imageFileName := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	// Consistent path structure as per project requirements
	saveDir := filepath.Join("attachments", subDir)
	imagePath := filepath.Join(saveDir, imageFileName)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}
	return filepath.Join(saveDir, imageFileName), nil
}

// formatToExtension maps an ImageFormat to a file extension.
func formatToExtension(format utils.ImageFormat) (string, bool) {
	switch format {
	case utils.JPEG:
		return ".jpg", true
	case utils.PNG:
		return ".png", true
	case utils.GIF:
		return ".gif", true
	case utils.WebP:
		return ".webp", true
	case utils.BMP:
		return ".bmp", true
	case utils.TIFF:
		return ".tiff", true
	default:
		return "", false
	}
}
