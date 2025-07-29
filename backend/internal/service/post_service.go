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

func (s *PostService) CreatePostWithViewers(post *models.Post, imageData []byte, imageMimeType string, viewerIDs []int64) (int64, error) {
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

	// Create the post
	postID, err := s.PostStore.CreatePost(post)
	if err != nil {
		return 0, err
	}

	// If it's a private post and has viewers, add them to Post_Visibility
	if post.Privacy == "private" && len(viewerIDs) > 0 {
		err = s.PostStore.AddPostViewers(postID, viewerIDs)
		if err != nil {
			return 0, fmt.Errorf("failed to add post viewers: %w", err)
		}
	}

	return postID, nil
}

func (s *PostService) SearchUsers(query string, currentUserID int64) ([]*models.User, error) {
	if query == "" {
		return []*models.User{}, nil
	}
	return s.PostStore.SearchUsers(query, currentUserID)
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

func (s *PostService) UpdatePost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Post, error) {
	// Get the existing post
	post, err := s.PostStore.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check if the user is authorized to edit this post
	if post.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate content
	if content == "" {
		return nil, fmt.Errorf("post content is required")
	}

	// Handle image update if provided
	var imagePath string
	if len(imageData) > 0 {
		savedImagePath, err := s.saveImage(imageData, "posts")
		if err != nil {
			return nil, err
		}
		imagePath = savedImagePath
	} else {
		// Keep existing image if no new image provided
		imagePath = post.Image
	}

	// Update the post in the store
	updatedPost, err := s.PostStore.UpdatePost(postID, content, imagePath)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
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
	saveDir := filepath.Join(subDir)
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
