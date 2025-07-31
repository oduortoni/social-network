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

func (ps *ProfileStore) MyProfileDetails(userid int64, profile models.ProfileDetails) (models.ProfileDetails, error) {
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

// Helper function to handle null string values
func getStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
