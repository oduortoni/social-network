package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// User represents a user in the database.
type User struct {
	ID              int64      `json:"id"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	FirstName       *string    `json:"first_name,omitempty"`
	LastName        *string    `json:"last_name,omitempty"`
	DateOfBirth     *string    `json:"date_of_birth,omitempty"`
	Avatar          *string    `json:"avatar,omitempty"`
	Nickname        *string    `json:"nickname,omitempty"`
	AboutMe         *string    `json:"about_me,omitempty"`
	IsProfilePublic bool       `json:"is_profile_public"`
	CreatedAt       time.Time  `json:"created_at"`
}

// String returns a pretty-printed multiline JSON-style representation of the User.
func (user User) String() string {
	data, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		return fmt.Sprintf("User{error: %v}", err)
	}
	return string(data)
}
