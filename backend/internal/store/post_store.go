package store

import (
	"database/sql"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
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
	row := s.DB.QueryRow("SELECT id, user_id, content, image, privacy, created_at, updated_at FROM Posts WHERE id = ?", id)

	var post models.Post
	var updatedAt sql.NullTime
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// Set the updated_at field and is_edited flag
	if updatedAt.Valid {
		post.UpdatedAt = &updatedAt.Time
		post.IsEdited = true
	}

	return &post, nil
}

func (s *PostStore) UpdatePost(postID int64, content, imagePath string) (*models.Post, error) {
	// Update the post with new content, image, and set updated_at timestamp
	stmt, err := s.DB.Prepare("UPDATE Posts SET content = ?, image = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	now := time.Now()
	_, err = stmt.Exec(content, imagePath, now, postID)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated post with author information
	row := s.DB.QueryRow(`
        SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at, p.updated_at,
               u.first_name, u.last_name, u.nickname, u.avatar
        FROM Posts p
        JOIN Users u ON p.user_id = u.id
        WHERE p.id = ?
    `, postID)

	var post models.Post
	var updatedAt sql.NullTime
	err = row.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy,
		&post.CreatedAt, &updatedAt, &post.Author.FirstName, &post.Author.LastName,
		&post.Author.Nickname, &post.Author.Avatar)
	if err != nil {
		return nil, err
	}

	// Set the updated_at field and is_edited flag
	if updatedAt.Valid {
		post.UpdatedAt = &updatedAt.Time
		post.IsEdited = true
	}

	return &post, nil
}

func (s *PostStore) GetPosts(userID int64) ([]*models.Post, error) {
	rows, err := s.DB.Query(`
        SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at, p.updated_at,
               u.first_name, u.last_name, u.nickname, u.avatar
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
		var updatedAt sql.NullTime
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy,
			&post.CreatedAt, &updatedAt, &post.Author.FirstName, &post.Author.LastName,
			&post.Author.Nickname, &post.Author.Avatar); err != nil {
			return nil, err
		}

		// Set the updated_at field and is_edited flag
		if updatedAt.Valid {
			post.UpdatedAt = &updatedAt.Time
			post.IsEdited = true
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *PostStore) GetCommentsByPostID(postID int64) ([]*models.Comment, error) {
	rows, err := s.DB.Query(`
        SELECT c.id, c.post_id, c.user_id, c.content, c.image, c.created_at, c.updated_at,
               u.first_name, u.last_name, u.nickname, u.avatar
        FROM Comments c
        JOIN Users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		var updatedAt sql.NullTime
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image,
			&comment.CreatedAt, &updatedAt, &comment.Author.FirstName, &comment.Author.LastName,
			&comment.Author.Nickname, &comment.Author.Avatar); err != nil {
			return nil, err
		}

		// Set the updated_at field and is_edited flag
		if updatedAt.Valid {
			comment.UpdatedAt = &updatedAt.Time
			comment.IsEdited = true
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (s *PostStore) DeletePost(postID int64) error {
	_, err := s.DB.Exec("DELETE FROM Posts WHERE id = ?", postID)
	return err
}

// AddPostViewers adds viewers to a private post
func (s *PostStore) AddPostViewers(postID int64, viewerIDs []int64) error {
	if len(viewerIDs) == 0 {
		return nil
	}

	// Prepare the insert statement
	stmt, err := s.DB.Prepare("INSERT INTO Post_Visibility (post_id, viewer_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert each viewer
	for _, viewerID := range viewerIDs {
		_, err := stmt.Exec(postID, viewerID)
		if err != nil {
			return err
		}
	}

	return nil
}

// SearchUsers searches for users by name or nickname
func (s *PostStore) SearchUsers(query string, currentUserID int64) ([]*models.User, error) {
	searchQuery := "%" + query + "%"

	rows, err := s.DB.Query(`
		SELECT id, first_name, last_name, nickname, avatar
		FROM Users
		WHERE id != ? AND (
			first_name LIKE ? OR
			last_name LIKE ? OR
			nickname LIKE ? OR
			(first_name || ' ' || last_name) LIKE ?
		)
		ORDER BY
			CASE
				WHEN nickname LIKE ? THEN 1
				WHEN first_name LIKE ? THEN 2
				WHEN last_name LIKE ? THEN 3
				ELSE 4
			END,
			first_name, last_name
		LIMIT 10
	`, currentUserID, searchQuery, searchQuery, searchQuery, searchQuery, searchQuery, searchQuery, searchQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Avatar)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// UpdateComment updates a comment's content and image, setting the updated_at timestamp
func (s *PostStore) UpdateComment(commentID int64, content, imagePath string) (*models.Comment, error) {
	// Update the comment with new content, image, and set updated_at timestamp
	stmt, err := s.DB.Prepare("UPDATE Comments SET content = ?, image = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	now := time.Now()
	_, err = stmt.Exec(content, imagePath, now, commentID)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated comment with author information
	row := s.DB.QueryRow(`
        SELECT c.id, c.post_id, c.user_id, c.content, c.image, c.created_at, c.updated_at,
               u.first_name, u.last_name, u.nickname, u.avatar
        FROM Comments c
        JOIN Users u ON c.user_id = u.id
        WHERE c.id = ?
    `, commentID)

	var comment models.Comment
	var updatedAt sql.NullTime
	err = row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image,
		&comment.CreatedAt, &updatedAt, &comment.Author.FirstName, &comment.Author.LastName,
		&comment.Author.Nickname, &comment.Author.Avatar)
	if err != nil {
		return nil, err
	}

	// Set the updated_at field and is_edited flag
	if updatedAt.Valid {
		comment.UpdatedAt = &updatedAt.Time
		comment.IsEdited = true
	}

	return &comment, nil
}

// DeleteComment removes a comment from the database
func (s *PostStore) DeleteComment(commentID int64) error {
	_, err := s.DB.Exec("DELETE FROM Comments WHERE id = ?", commentID)
	return err
}

// GetCommentByID retrieves a specific comment by its ID with author information
func (s *PostStore) GetCommentByID(commentID int64) (*models.Comment, error) {
	row := s.DB.QueryRow(`
        SELECT c.id, c.post_id, c.user_id, c.content, c.image, c.created_at, c.updated_at,
               u.first_name, u.last_name, u.nickname, u.avatar
        FROM Comments c
        JOIN Users u ON c.user_id = u.id
        WHERE c.id = ?
    `, commentID)

	var comment models.Comment
	var updatedAt sql.NullTime
	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image,
		&comment.CreatedAt, &updatedAt, &comment.Author.FirstName, &comment.Author.LastName,
		&comment.Author.Nickname, &comment.Author.Avatar)
	if err != nil {
		return nil, err
	}

	// Set the updated_at field and is_edited flag
	if updatedAt.Valid {
		comment.UpdatedAt = &updatedAt.Time
		comment.IsEdited = true
	}

	return &comment, nil
}
