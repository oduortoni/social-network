package store

import "github.com/tajjjjr/social-network/backend/internal/models"

// PostStoreInterface defines the interface for post-related database operations.
type PostStoreInterface interface {
	CreatePost(post *models.Post) (int64, error)
	CreateComment(comment *models.Comment) (int64, error)
	GetPostByID(id int64) (*models.Post, error)
	GetPosts(userID int64) ([]*models.Post, error)
	UpdatePost(postID int64, content, imagePath string) (*models.Post, error)
	GetCommentsByPostID(postID int64) ([]*models.Comment, error)
	DeletePost(postID int64) error
	AddPostViewers(postID int64, viewerIDs []int64) error
	SearchUsers(query string, currentUserID int64) ([]*models.User, error)
}
