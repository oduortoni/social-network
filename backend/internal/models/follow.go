package models

type Followee struct {
	FolloweeId int `json:"followeeid"`
}

type FollowFollowingStat struct {
	NumberOfFollowers int `json:"numberoffollowers"`
	NumberOfFollowing int `json:"numberoffollowing"`
}

type FollowListUser struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Avatar     string `json:"avatar"`
	FollowerID int64  `json:"follower_id"`
}

type FollowListResponse struct {
	Followers []FollowListUser `json:"user"`
}
