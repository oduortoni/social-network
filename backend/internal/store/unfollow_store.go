package store

import "database/sql"

// UnfollowStore handles database operations for unfollow.
type UnfollowStore struct {
	DB *sql.DB
}

// NewUnfollowStore creates a new UnfollowStore.
func NewUnfollowStore(db *sql.DB) *UnfollowStore {
	return &UnfollowStore{DB: db}
}

// GetFollowConnectionID retrieves the follow connection ID between two users
func (unf *UnfollowStore) GetFollowConnectionID(followerID, followeeID int64) (int64, error) {
	var connectionID int64
	query := "SELECT id FROM Followers WHERE follower_id = ? AND followee_id = ?"
	err := unf.DB.QueryRow(query, followerID, followeeID).Scan(&connectionID)
	if err != nil {
		return 0, err
	}
	return connectionID, nil
}

// DeleteFollowConnection removes a follow connection from the database
func (unf *UnfollowStore) DeleteFollowConnection(followConnectionID int64) error {
	query := "DELETE FROM Followers WHERE id = ?"
	_, err := unf.DB.Exec(query, followConnectionID)
	return err
}
