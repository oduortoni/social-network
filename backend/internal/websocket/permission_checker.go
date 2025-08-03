package websocket

import "database/sql"

type DBPermissionChecker struct {
	DB *sql.DB
}

func NewDBPermissionChecker(db *sql.DB) *DBPermissionChecker {
	return &DBPermissionChecker{DB: db}
}

// CanUsersChat checks if two users are allowed to chat.
// This means they must be followers of each other OR the target has a public profile.
func (p *DBPermissionChecker) CanUsersChat(userID, targetID int64) (bool, error) {
	if userID == targetID {
		return true, nil // Users can always chat with themselves (e.g., for notes)
	}

	// Check if they are mutual followers
	var mutualFollowCount int
	err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM Followers
		WHERE (follower_id = ? AND followee_id = ? AND status = 'accepted')
		   OR (follower_id = ? AND followee_id = ? AND status = 'accepted')
	`, userID, targetID, targetID, userID).Scan(&mutualFollowCount)
	if err != nil {
		return false, err
	}

	if mutualFollowCount == 2 {
		return true, nil // Both are following each other
	}

	// Check if target user has a public profile
	var isPublic bool
	err = p.DB.QueryRow(`
		SELECT is_profile_public FROM Users WHERE id = ?
	`, targetID).Scan(&isPublic)
	if err != nil {
		return false, err
	}

	return isPublic, nil // Allow chat if target profile is public
}
