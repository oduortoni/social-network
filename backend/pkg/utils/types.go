package utils

type UserID string

type Response struct {
	Message string `json:"message,omitempty"`
}

var User_id UserID = "userID"
