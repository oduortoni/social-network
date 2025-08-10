package store

import (
	"database/sql"
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type groupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) GroupStore {
	return &groupStore{db: db}
}

func (s *groupStore) CreateGroup(group *models.Group) (*models.Group, error) {
	stmt, err := s.db.Prepare("INSERT INTO groups (creator_id, title, description, privacy) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(group.CreatorID, group.Title, group.Description, group.Privacy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	group.ID = int(id)

	return group, nil
}

func (s *groupStore) GetGroupByID(groupID int) (*models.Group, error) {
	var group models.Group
	err := s.db.QueryRow("SELECT id, creator_id, title, description, privacy, created_at FROM groups WHERE id = ?", groupID).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.Privacy,
		&group.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, err
	}
	return &group, nil
}
