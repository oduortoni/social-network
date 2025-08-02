package websocket

import "database/sql"

type DBPermissionChecker struct {
	DB *sql.DB
}

func NewDBPermissionChecker(db *sql.DB) *DBPermissionChecker {
	return &DBPermissionChecker{DB: db}
}

// CanUsersChat checks if two users are allowed to chat.
// This means they must be followers of each other.
func (p *DBPermissionChecker) CanUsersChat(userID, targetID int64) (bool, error) {
	if userID == targetID {
		return true, nil // Users can always chat with themselves (e.g., for notes)
	}

	var count int
	err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM Followers
		WHERE (follower_id = ? AND followee_id = ? AND status = 'accepted')
		   OR (follower_id = ? AND followee_id = ? AND status = 'accepted')
	`, userID, targetID, targetID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 2, nil // Both must be following each other
}
