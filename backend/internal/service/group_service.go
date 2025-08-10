package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type groupService struct {
	groupStore store.GroupStore
}

func NewGroupService(groupStore store.GroupStore) GroupService {
	return &groupService{groupStore: groupStore}
}

func (s *groupService) CreateGroup(group *models.Group) (*models.Group, error) {
	return s.groupStore.CreateGroup(group)
}
