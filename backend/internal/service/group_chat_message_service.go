package service

import (
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type groupChatMessageService struct {
	groupChatMessageStore store.GroupChatMessageStore
	groupService          GroupService
	groupMemberStore      store.GroupMemberStore
}

func NewGroupChatMessageService(groupChatMessageStore store.GroupChatMessageStore, groupService GroupService, groupMemberStore store.GroupMemberStore) GroupChatMessageService {
	return &groupChatMessageService{groupChatMessageStore: groupChatMessageStore, groupService: groupService, groupMemberStore: groupMemberStore}
}

func (s *groupChatMessageService) SendGroupChatMessage(groupID, senderID int64, content string) (*models.GroupChatMessage, error) {
	_, err := s.groupService.GetGroupByID(int64(groupID))
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	isMember, err := s.groupMemberStore.IsGroupMember(int64(groupID), int64(senderID))
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	message := &models.GroupChatMessage{
		GroupID:  groupID,
		SenderID: senderID,
		Content:  content,
	}

	createdMessage, err := s.groupChatMessageStore.CreateGroupChatMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send group chat message: %w", err)
	}

	return createdMessage, nil
}

func (s *groupChatMessageService) GetGroupChatMessages(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error) {
	_, err := s.groupService.GetGroupByID(int64(groupID))
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	isMember, err := s.groupMemberStore.IsGroupMember(groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	messages, err := s.groupChatMessageStore.GetGroupChatMessages(groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get group chat messages: %w", err)
	}

	return messages, nil
}
