package store

import (
	"database/sql"
	"strings"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
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
	_, err := followstore.DB.Exec("INSERT INTO Followers (follower_id, followee_id, status, requested_at, accepted_at) VALUES (?, ?, ?, ?, ?)", followerId, followeeId, "accepted", currentTime, currentTime)
	return err
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
	var firstName, lastName sql.NullString
	query := "SELECT first_name, last_name FROM Users WHERE id = ?"
	err := followstore.DB.QueryRow(query, userID).Scan(&firstName, &lastName)
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

	return name, name, nil
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


func (followstore *FollowStore) GetUserFollowers(userid int64) (models.FollowListResponse, error) {
	var followersList models.FollowListResponse
	rows, err := followstore.DB.Query(`
		SELECT u.id, u.first_name, u.last_name, u.avatar 
		FROM Users u 
		INNER JOIN Followers f ON u.id = f.follower_id 
		WHERE f.followee_id = ? AND f.status = 'accepted'`, userid)
	if err != nil {
		return followersList, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower models.FollowUser
		var firstName, lastName, avatar sql.NullString
		err := rows.Scan(&follower.FollowerID, &firstName, &lastName, &avatar)
		if err != nil {
			return followersList, err
		}

		if firstName.Valid {
			follower.FirstName = firstName.String
		}
		if lastName.Valid {
			follower.LastName = lastName.String
		}
		if avatar.Valid {
			follower.Avatar = avatar.String
		}

		followersList.Followers = append(followersList.Followers, follower)
	}

	if err = rows.Err(); err != nil {
		return followersList, err
	}

	return followersList, nil
}


func (followstore *FollowStore) GetUserFollowees(userid int64) (models.FollowListResponse, error) {
	var followersList models.FollowListResponse
	rows, err := followstore.DB.Query(`
		SELECT u.id, u.first_name, u.last_name, u.avatar 
		FROM Users u 
		INNER JOIN Followers f ON u.id = f.followee_id 
		WHERE f.follower_id = ? AND f.status = 'accepted'`, userid)
	if err != nil {
		return followersList, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower models.FollowUser
		var firstName, lastName, avatar sql.NullString
		err := rows.Scan(&follower.FollowerID, &firstName, &lastName, &avatar)
		if err != nil {
			return followersList, err
		}

		if firstName.Valid {
			follower.FirstName = firstName.String
		}
		if lastName.Valid {
			follower.LastName = lastName.String
		}
		if avatar.Valid {
			follower.Avatar = avatar.String
		}

		followersList.Followers = append(followersList.Followers, follower)
	}

	if err = rows.Err(); err != nil {
		return followersList, err
	}

	return followersList, nil
}
