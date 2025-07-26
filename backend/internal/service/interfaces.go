package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
)

// AuthServiceInterface defines the interface for the auth service.
type AuthServiceInterface interface {
	AuthenticateUser(email, password string) (*models.User, string, error)
	DeleteSession(sessionID string) (int, error)
	GetUserIDBySession(sessionID string) (int, error)
	CreateUser(user *models.User) (*models.User, error)
	ValidateEmail(email string) (bool, error)
	UserExists(email string) (bool, error)
}

// PostServiceInterface defines the interface for the post service.
type PostServiceInterface interface {
	CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error)
	CreatePostWithViewers(post *models.Post, imageData []byte, imageMimeType string, viewerIDs []int64) (int64, error)
	GetPostByID(id int64) (*models.Post, error)
	GetPosts(userID int64) ([]*models.Post, error)
	CreateComment(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error)
	GetCommentsByPostID(postID int64) ([]*models.Comment, error)
	DeletePost(postID, userID int64) error
	SearchUsers(query string, currentUserID int64) ([]*models.User, error)
}

type FollowServiceInterface interface {
	IsAccountPublic(followeeID int64) (bool, error)
	CreateFollowForPublicAccount(followerid, followeeid int64) error
	CreateFollowForPrivateAccount(followrid, followeeid int64) (int64, error)
}

type UnfollowServiceInterface interface {
	GetFollowConnectionID(followerID, followeeID int64) (int64, error)
	DeleteFollowConnection(followConnectionID int64) error
}

type FollowRequestServiceInterface interface {
	AcceptedFollowConnection(followConnectionID int64) error
	RejectedFollowConnection(followConnectionID int64) error
}
