package models

type Followee struct {
	FolloweeId int `json:"followeeid"`
}

type FollowFollowingStat struct {
	NumberOfFollowers int `json:"numberoffollowers"`
	NumberOfFollowing int `json:"numberoffollowing"`
}
