package store

import (
	"database/sql"
	"time"
)

// FollowRequestStore handles database operations for follow request.
type FollowRequestStore struct {
	DB *sql.DB
}

// NewFollowRequestStore creates a new FollowRequestStore.
func NewFollowRequestStore(db *sql.DB) *FollowRequestStore {
	return &FollowRequestStore{DB: db}
}

// AcceptFollowConnection accepts a follow request by setting status to 'accepted' and updating accepted_at
func (fr *FollowRequestStore) AcceptFollowConnection(followConnectionID int64) error {
	query := "UPDATE Followers SET status = 'accepted', accepted_at = ? WHERE id = ? AND status = 'pending'"
	_, err := fr.DB.Exec(query, time.Now(), followConnectionID)
	return err
}

// RejectFollowConnection rejects a follow request by setting status to 'rejected'
func (fr *FollowRequestStore) RejectFollowConnection(followConnectionID int64) error {
	query := "UPDATE Followers SET status = 'rejected' WHERE id = ? AND status = 'pending'"
	_, err := fr.DB.Exec(query, followConnectionID)
	return err
}

// GetPendingFollowRequests retrieves all pending follow requests for a user
func (fr *FollowRequestStore) GetPendingFollowRequests(followeeID int64) ([]int64, error) {
	query := "SELECT id FROM Followers WHERE followee_id = ? AND status = 'pending'"
	rows, err := fr.DB.Query(query, followeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requestIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		requestIDs = append(requestIDs, id)
	}

	return requestIDs, rows.Err()
}
