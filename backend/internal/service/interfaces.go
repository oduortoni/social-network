package service

import (
	

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// AuthServiceInterface defines the interface for the auth service.
type AuthServiceInterface interface {
	AuthenticateUser(email, password string) (*models.User, string, error)
}

// PostServiceInterface defines the interface for the post service.
type PostServiceInterface interface {
	CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error)
	GetPostByID(id int64) (*models.Post, error)
	GetFeed(userID int64) ([]*models.Post, error)
}
