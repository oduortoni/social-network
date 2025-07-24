package store

import (
	"database/sql"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"time"
)

type PostStore struct {
	DB *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{DB: db}
}

func (s *PostStore) CreatePost(post *models.Post) (int64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO Posts (user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.UserID, post.Content, post.Image, post.Privacy, time.Now())
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *PostStore) CreateComment(comment *models.Comment) (int64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO Comments (post_id, user_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(comment.PostID, comment.UserID, comment.Content, comment.Image, time.Now())
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *PostStore) GetPostByID(id int64) (*models.Post, error) {
	row := s.DB.QueryRow("SELECT id, user_id, content, image, privacy, created_at FROM Posts WHERE id = ?", id)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostStore) GetFeed(userID int64) ([]*models.Post, error) {
	rows, err := s.DB.Query(`
		SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at
		FROM Posts p
		INNER JOIN Followers f ON p.user_id = f.followee_id
		WHERE f.follower_id = ? AND f.is_accepted = 1
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
