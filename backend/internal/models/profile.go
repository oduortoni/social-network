package models

type ProfileDetails struct {
	FirstName string `json:"firstname"` // User's first name
	LastName  string `json:"lastname"`  // User's last name
	Email     string `json:"email"`
	Avatar    string `json:"avatar"` // URL or path to user's avatar image
	ID        int64  `json:"id"`     // User's unique ID
	About     string `json:"about"`  // User's bio or description
	Nickname  string `json:"nickname"`
	// Follow button status:
	// - "hide"     → Hidden (for the logged-in user's own profile)
	// - "pending"  → Follow request sent but not yet accepted
	// - "follow"   → Not following, follow option available
	// - "following"→ Already following the user
	FollowbtnStatus string `json:"followbtnstatus"`
	// Message button visibility:
	// - "hide"     → Hidden (for the logged-in user's own profile)
	// - "visible"  → Visible to followers (messaging allowed)
	MessageBtnStatus  string `json:"messagebtnstatus"`
	DateOfBirth       string `json:"dateofbirth"`
	Profile           string `json:"profile"`
	NumberOfFollowers int    `json:"numberoffollowers"`
	NumberOfFollowees int    `json:"numberoffollowees"`
	NumberOfPosts     int    `json:"numberofposts"`
}

type ProfileResponse struct {
	ProfileDetails ProfileDetails `json:"profile_details"`
	// Posts          []Post         `json:"posts"`
}
