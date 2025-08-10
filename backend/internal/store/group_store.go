package store

import (
	"database/sql"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) GroupStore {
	return &groupStore{db: db}
}

func (s *groupStore) CreateGroup(group *models.Group) (*models.Group, error) {
	stmt, err := s.db.Prepare("INSERT INTO groups (creator_id, title, description) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(group.CreatorID, group.Title, group.Description)
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
