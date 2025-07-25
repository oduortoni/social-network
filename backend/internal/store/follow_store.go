package store

import (
	"database/sql"
	"time"
)

// FollowStore handles database operations for follow and unfollow.
type FollowStore struct {
	DB *sql.DB
}

// NewFollowStore creates a new FollowStore.
func NewFollowStore(db *sql.DB) *FollowStore {
	return &FollowStore{DB: db}
}

func (followstore *FollowStore) IsUserAccountPublic(userid int64) (bool, error) {
	var num int
	err := followstore.DB.QueryRow("SELECT is_profile_public FROM Users WHERE id=?", userid).Scan(&num)
	if err != nil {
		return false, err
	}
	return num == 1, nil
}

func (followstore *FollowStore) CreatePublicFollowConnection(followerId, followeeId int64) error {
	currentTime := time.Now()
	_, err := followstore.DB.Exec("INSERT INTO Followers (follower_id, followee_id,is_accepted,requested_at,accepted_at) VALUES (?, ?,?,?,?)", followerId, followeeId, 1, currentTime, currentTime)
	return err
}
