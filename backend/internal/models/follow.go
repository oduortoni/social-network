package models

type Followee struct {
	FolloweeId int `json:"followeeid"`
}

type FollowUser struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Avatar     string `json:"avatar"`
	FollowerID int64  `json:"follower_id"`
}

type FollowListResponse struct {
	Followers []FollowUser `json:"user"`
}
