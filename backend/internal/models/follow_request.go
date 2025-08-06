package models

type FollowRequestResponseStatus struct {
	Status string `json:"status"`
}

type FollowUserRequest struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Avatar      string `json:"avatar"`
	FollowerID  int64  `json:"follower_id"`
	RequestedAt string `json:"requested_at"`
	RequestID   int64  `json:"request_id"`
}

type FollowRequestUserResponse struct {
	Followers []FollowUserRequest `json:"user"`
}
