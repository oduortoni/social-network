package store

import (
	"database/sql"
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// ProfileStore handles database operations for profile.
type ProfileStore struct {
	DB *sql.DB
}

// NewProfileStore creates a new ProfileStore.
func NewProfileStore(db *sql.DB) *ProfileStore {
	return &ProfileStore{DB: db}
}

func (ps *ProfileStore) MyProfileDetails(userid int64) (models.ProfileDetails, error) {
	// Initialize an empty ProfileDetails struct
	var profile models.ProfileDetails
	// Prepare the SQL query to fetch user details
	var firstName, lastName, email, nickname, aboutMe, avatar sql.NullString
	var dateOfBirth sql.NullTime
	var isProfilePublic bool

	query := `SELECT first_name, last_name, email, nickname, about_me, date_of_birth, is_profile_public, avatar 
			  FROM Users 
			  WHERE id = ?`

	err := ps.DB.QueryRow(query, userid).Scan(
		&firstName,
		&lastName,
		&email,
		&nickname,
		&aboutMe,
		&dateOfBirth,
		&isProfilePublic,
		&avatar,
	)
	if err != nil {
		return profile, err
	}

	// Convert the retrieved data to strings, handling null values
	profile.FirstName = getStringValue(firstName)
	profile.LastName = getStringValue(lastName)
	profile.Email = getStringValue(email)
	profile.Nickname = getStringValue(nickname)
	profile.About = getStringValue(aboutMe)
	profile.DateOfBirth = dateOfBirth.Time.Format("2006-01-02")
	profile.Profile = fmt.Sprintf("%v", isProfilePublic)
	profile.Avatar = getStringValue(avatar)

	return profile, nil
}

func (followStore *ProfileStore) GetFollowersStat(userid int64) (int, int, error) {
	followers := 0
	following := 0
	err := followStore.DB.QueryRow("SELECT COUNT(*) FROM Followers WHERE follower_id = ? AND status = 'accepted'", userid).Scan(&followers)
	if err != nil {
		return 0, 0, err
	}
	err = followStore.DB.QueryRow("SELECT COUNT(*) FROM Followers WHERE followee_id = ? AND status = 'accepted'", userid).Scan(&following)
	if err != nil {
		return 0, 0, err
	}
	return followers, following, nil
}

func (ps *ProfileStore) GetNumberOfPosts(userid int64) (int, error) {
	var count int
	err := ps.DB.QueryRow("SELECT COUNT(*) FROM Posts WHERE user_id = ?", userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ps *ProfileStore) GetFollowStatus(user_id, LoggedInUser int64) (string, error) {
	var status string
	query := `SELECT status FROM Followers WHERE follower_id = ? AND followee_id = ?`
	err := ps.DB.QueryRow(query, LoggedInUser, user_id).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "follow", nil // No follow relationship found
		}
		return "", err // Other error
	}

	if status == "accepted" {
		return "following", nil // Already following
	} else if status == "pending" {
		return "pending", nil // Follow request is pending
	}
	return "follow", nil // Follow request was rejected, can follow again
}

// Helper function to handle null string values
func getStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
