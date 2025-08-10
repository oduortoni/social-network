package service

import (
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type groupRequestService struct {
	groupRequestStore store.GroupRequestStore
	groupService      GroupService
}

func NewGroupRequestService(groupRequestStore store.GroupRequestStore, groupService GroupService) GroupRequestService {
	return &groupRequestService{groupRequestStore: groupRequestStore, groupService: groupService}
}

func (s *groupRequestService) SendJoinRequest(groupID, userID int) (*models.GroupRequest, error) {
	group, err := s.groupService.GetGroupByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	if group.Privacy == "private" {
		return nil, fmt.Errorf("cannot send join request to a private group")
	}

	// TODO: Add logic to check if the user is already a member or has a pending request.

	request := &models.GroupRequest{
		GroupID: groupID,
		UserID:  userID,
		Status:  "pending",
	}

	createdRequest, err := s.groupRequestStore.CreateGroupRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create group request: %w", err)
	}

	return createdRequest, nil
}

func (s *groupRequestService) ApproveJoinRequest(requestID int, approverID int) error {
	// TODO: Add logic to verify approverID is an admin/creator of the group.
	request, err := s.groupRequestStore.GetGroupRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to get group request: %w", err)
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not pending")
	}

	// TODO: Add the user to the group_members table.

	err = s.groupRequestStore.UpdateGroupRequestStatus(requestID, "approved")
	if err != nil {
		return fmt.Errorf("failed to approve group request: %w", err)
	}

	return nil
}

func (s *groupRequestService) RejectJoinRequest(requestID int, rejecterID int) error {
	// TODO: Add logic to verify rejecterID is an admin/creator of the group.
	request, err := s.groupRequestStore.GetGroupRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to get group request: %w", err)
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not pending")
	}

	err = s.groupRequestStore.UpdateGroupRequestStatus(requestID, "rejected")
	if err != nil {
		return fmt.Errorf("failed to reject group request: %w", err)
	}

	return nil
}
