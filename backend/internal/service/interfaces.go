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
	CreatePost(post *models.Post) (int64, error)
}
