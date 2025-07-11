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
	stmt, err := s.DB.Prepare("INSERT INTO POSTS (user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)")
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

func (s *PostStore) GetPostByID(id int64) (*models.Post, error) {
	row := s.DB.QueryRow("SELECT id, user_id, content, image, privacy, created_at FROM POSTS WHERE id = ?", id)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
