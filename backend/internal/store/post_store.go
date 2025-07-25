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

func (s *PostStore) GetPosts(userID int64) ([]*models.Post, error) {
	rows, err := s.DB.Query(`
        SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at, u.nickname, u.avatar
        FROM Posts p
        JOIN Users u ON p.user_id = u.id
        WHERE p.privacy = 'public'
        OR (p.privacy = 'almost_private' AND p.user_id = ? OR p.user_id IN (
            SELECT follower_id FROM Followers WHERE followee_id = ? AND is_accepted = 1
        ))
        OR (p.privacy = 'private' AND EXISTS (
            SELECT 1 FROM Post_visibility pv WHERE pv.post_id = p.id AND pv.viewer_id = ?
        ))
        ORDER BY p.created_at DESC
    `, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &post.Author.Nickname, &post.Author.Avatar); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *PostStore) GetCommentsByPostID(postID int64) ([]*models.Comment, error) {
    rows, err := s.DB.Query(`
        SELECT c.id, c.post_id, c.user_id, c.content, c.image, c.created_at, u.nickname, u.avatar
        FROM Comments c
        JOIN Users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at ASC
    `, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []*models.Comment
    for rows.Next() {
        var comment models.Comment
        if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image, &comment.CreatedAt, &comment.Author.Nickname, &comment.Author.Avatar); err != nil {
            return nil, err
        }
        comments = append(comments, &comment)
    }

    return comments, nil
}

func (s *PostStore) DeletePost(postID int64) error {
	_, err := s.DB.Exec("DELETE FROM Posts WHERE id = ?", postID)
	return err
}

