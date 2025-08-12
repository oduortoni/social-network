package store

import (
	"database/sql"
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupMemberStore struct {
	db *sql.DB
}

func NewGroupMemberStore(db *sql.DB) GroupMemberStore {
	return &groupMemberStore{db: db}
}

func (s *groupMemberStore) IsGroupMember(groupID, userID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)", groupID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking group membership: %w", err)
	}
	return exists, nil
}

func (s *groupMemberStore) AddGroupMember(groupID, userID int64, role string) (*models.GroupMember, error) {
	stmt, err := s.db.Prepare("INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(groupID, userID, role)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	member := &models.GroupMember{
		ID:      id,
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}
	return member, nil
}

func (s *groupMemberStore) RemoveGroupMember(groupID, userID int64) error {
	stmt, err := s.db.Prepare("DELETE FROM group_members WHERE group_id = ? AND user_id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(groupID, userID)
	if err != nil {
		return fmt.Errorf("error removing group member: %w", err)
	}
	return nil
}
