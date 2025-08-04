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
}

type FollowRequestUserResponse struct {
	Followers []FollowUserRequest `json:"user"`
}
