package store

import (
	"database/sql"
	"strings"
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

func (followstore *FollowStore) CreatePublicFollowConnection(followerId, followeeId int64)(int64, error) {
	currentTime := time.Now()
	result, err := followstore.DB.Exec("INSERT INTO Followers (follower_id, followee_id, status, requested_at, accepted_at) VALUES (?, ?, ?, ?, ?)", followerId, followeeId, "accepted", currentTime, currentTime)
     if err != nil {
		return 0, err
	}

	followID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return followID, nil
}

func (followstore *FollowStore) CreatePrivateFollowConnection(followerId, followeeId int64) (int64, error) {
	currentTime := time.Now()

	result, err := followstore.DB.Exec(
		"INSERT INTO Followers (follower_id, followee_id, status, requested_at) VALUES (?, ?, ?, ?)",
		followerId, followeeId, "pending", currentTime,
	)
	if err != nil {
		return 0, err
	}

	followID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return followID, nil
}

func (followstore *FollowStore) UserInfo(userID int64) (string, string, error) {
	var firstName, lastName, avatar sql.NullString
	query := "SELECT first_name, last_name, avatar FROM Users WHERE id = ?"
	err := followstore.DB.QueryRow(query, userID).Scan(&firstName, &lastName, &avatar)
	if err != nil {
		return "", "", err
	}

	name := ""
	if firstName.Valid && lastName.Valid {
		name = firstName.String + " " + lastName.String
	} else if firstName.Valid {
		name = firstName.String
	} else if lastName.Valid {
		name = lastName.String
	} else {
		name = "User"
	}

	return name, avatar.String, nil
}

func (followstore *FollowStore) AddtoNotification(follower_id int64, message string) error {
	// Determine notification type based on message content
	notificationType := "follow"
	if strings.Contains(message, "follow request") {
		notificationType = "follow_request"
	}

	_, err := followstore.DB.Exec(`INSERT INTO Notifications (user_id, type, message) VALUES (?, ?, ?)`,
		follower_id, notificationType, message)
	return err
}
