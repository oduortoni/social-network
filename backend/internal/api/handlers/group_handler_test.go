package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockGroupService is a mock implementation of the GroupService for testing.
type MockGroupService struct {
	CreateGroupFunc  func(group *models.Group) (*models.Group, error)
	GetGroupByIDFunc func(groupID int) (*models.Group, error)
}

func (m *MockGroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	if m.CreateGroupFunc != nil {
		return m.CreateGroupFunc(group)
	}
	return group, nil
}

func (m *MockGroupService) GetGroupByID(groupID int) (*models.Group, error) {
	if m.GetGroupByIDFunc != nil {
		return m.GetGroupByIDFunc(groupID)
	}
	return nil, errors.New("GetGroupByID not implemented")
}

// MockGroupRequestService is a mock implementation of the GroupRequestService for testing.
type MockGroupRequestService struct {
	SendJoinRequestFunc    func(groupID, userID int) (*models.GroupRequest, error)
	ApproveJoinRequestFunc func(requestID int, approverID int) error
	RejectJoinRequestFunc  func(requestID int, rejecterID int) error
}

func (m *MockGroupRequestService) SendJoinRequest(groupID, userID int) (*models.GroupRequest, error) {
	if m.SendJoinRequestFunc != nil {
		return m.SendJoinRequestFunc(groupID, userID)
	}
	return nil, errors.New("SendJoinRequest not implemented")
}

func (m *MockGroupRequestService) ApproveJoinRequest(requestID int, approverID int) error {
	if m.ApproveJoinRequestFunc != nil {
		return m.ApproveJoinRequestFunc(requestID, approverID)
	}
	return errors.New("ApproveJoinRequest not implemented")
}

func (m *MockGroupRequestService) RejectJoinRequest(requestID int, rejecterID int) error {
	if m.RejectJoinRequestFunc != nil {
		return m.RejectJoinRequestFunc(requestID, rejecterID)
	}
	return errors.New("RejectJoinRequest not implemented")
}

// MockGroupChatMessageService is a mock implementation of the GroupChatMessageService for testing.
type MockGroupChatMessageService struct {
	SendGroupChatMessageFunc func(groupID, senderID int, content string) (*models.GroupChatMessage, error)
	GetGroupChatMessagesFunc func(groupID int, userID int, limit, offset int) ([]*models.GroupChatMessage, error)
}

func (m *MockGroupChatMessageService) SendGroupChatMessage(groupID, senderID int, content string) (*models.GroupChatMessage, error) {
	if m.SendGroupChatMessageFunc != nil {
		return m.SendGroupChatMessageFunc(groupID, senderID, content)
	}
	return nil, errors.New("SendGroupChatMessage not implemented")
}

func (m *MockGroupChatMessageService) GetGroupChatMessages(groupID int, userID int, limit, offset int) ([]*models.GroupChatMessage, error) {
	if m.GetGroupChatMessagesFunc != nil {
		return m.GetGroupChatMessagesFunc(groupID, userID, limit, offset)
	}
	return nil, errors.New("GetGroupChatMessages not implemented")
}

func TestCreateGroup(t *testing.T) {
	// Test case 1: Successful group creation
	t.Run("Successful group creation", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			CreateGroupFunc: func(group *models.Group) (*models.Group, error) {
				group.ID = 1
				return group, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		group := &models.Group{
			Title:       "Test Group",
			Description: "This is a test group.",
			CreatorID:   1,
			Privacy:     "public",
		}

		body, _ := json.Marshal(group)
		req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Add user ID to context
		ctx := context.WithValue(req.Context(), "userID", 1)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var createdGroup models.Group
		if err := json.NewDecoder(rr.Body).Decode(&createdGroup); err != nil {
			t.Fatal(err)
		}

		if createdGroup.ID != 1 {
			t.Errorf("handler returned unexpected group ID: got %v want %v",
				createdGroup.ID, 1)
		}
	})
}

func TestSendJoinRequest(t *testing.T) {
	// Test case 1: Successful join request for a public group
	t.Run("Successful join request for public group", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "public"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{
			SendJoinRequestFunc: func(groupID, userID int) (*models.GroupRequest, error) {
				return &models.GroupRequest{ID: 1, GroupID: groupID, UserID: userID, Status: "pending"}, nil
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("POST /groups/{groupID}/join-request", http.HandlerFunc(h.SendJoinRequest))

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["id"] != float64(1) {
			t.Errorf("Expected request ID 1, got %v", response["id"])
		}
	})

	// Test case 2: Failed join request for a private group
	t.Run("Failed join request for private group", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{
			SendJoinRequestFunc: func(groupID, userID int) (*models.GroupRequest, error) {
				return nil, errors.New("cannot send join request to a private group")
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("POST /groups/{groupID}/join-request", http.HandlerFunc(h.SendJoinRequest))

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("cannot send join request to a private group")) {
			t.Errorf("Expected error message not found in response: %s", rr.Body.String())
		}
	})
}

func TestSendGroupChatMessage(t *testing.T) {
	// Test case 1: Successful message sending
	t.Run("Successful message sending", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			SendGroupChatMessageFunc: func(groupID, senderID int, content string) (*models.GroupChatMessage, error) {
				return &models.GroupChatMessage{ID: 1, GroupID: groupID, SenderID: senderID, Content: content}, nil
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("POST /groups/{groupID}/chat", http.HandlerFunc(h.SendGroupChatMessage))

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/1/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response models.GroupChatMessage
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Content != "Hello Group!" {
			t.Errorf("Expected message content 'Hello Group!', got %s", response.Content)
		}
	})

	// Test case 2: Failed message sending to non-existent group
	t.Run("Failed message sending to non-existent group", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return nil, errors.New("group not found")
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			SendGroupChatMessageFunc: func(groupID, senderID int, content string) (*models.GroupChatMessage, error) {
				return nil, errors.New("group not found")
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("POST /groups/{groupID}/chat", http.HandlerFunc(h.SendGroupChatMessage))

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/999/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("group not found")) {
			t.Errorf("Expected error message not found in response: %s", rr.Body.String())
		}
	})

	// Test case 3: Failed message sending by a non-member
	t.Run("Failed message sending by a non-member", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			SendGroupChatMessageFunc: func(groupID, senderID int, content string) (*models.GroupChatMessage, error) {
				return nil, errors.New("user is not a member of this group")
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("POST /groups/{groupID}/chat", http.HandlerFunc(h.SendGroupChatMessage))

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/1/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("user is not a member of this group")) {
			t.Errorf("Expected error message not found in response: %s", rr.Body.String())
		}
	})
}

func TestGetGroupChatMessages(t *testing.T) {
	// Test case 1: Successful message retrieval
	t.Run("Successful message retrieval", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			GetGroupChatMessagesFunc: func(groupID int, userID int, limit, offset int) ([]*models.GroupChatMessage, error) {
				return []*models.GroupChatMessage{
					{ID: 1, GroupID: groupID, SenderID: 101, Content: "Message 1"},
					{ID: 2, GroupID: groupID, SenderID: 102, Content: "Message 2"},
				}, nil
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("GET /groups/{groupID}/chat", http.HandlerFunc(h.GetGroupChatMessages))

		req, err := http.NewRequest("GET", "/groups/1/chat?limit=10&offset=0", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response []models.GroupChatMessage
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if len(response) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(response))
		}
		if response[0].Content != "Message 1" {
			t.Errorf("Expected first message content 'Message 1', got %s", response[0].Content)
		}
	})

	// Test case 2: Failed message retrieval from a non-existent group
	t.Run("Failed message retrieval from non-existent group", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return nil, errors.New("group not found")
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			GetGroupChatMessagesFunc: func(groupID int, userID int, limit, offset int) ([]*models.GroupChatMessage, error) {
				return nil, errors.New("group not found")
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("GET /groups/{groupID}/chat", http.HandlerFunc(h.GetGroupChatMessages))

		req, err := http.NewRequest("GET", "/groups/999/chat", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("group not found")) {
			t.Errorf("Expected error message not found in response: %s", rr.Body.String())
		}
	})

	// Test case 3: Failed message retrieval by a non-member
	t.Run("Failed message retrieval by a non-member", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			GetGroupChatMessagesFunc: func(groupID int, userID int, limit, offset int) ([]*models.GroupChatMessage, error) {
				return nil, errors.New("user is not a member of this group")
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		mux := http.NewServeMux()
		mux.Handle("GET /groups/{groupID}/chat", http.HandlerFunc(h.GetGroupChatMessages))

		req, err := http.NewRequest("GET", "/groups/1/chat", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), "userID", 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("user is not a member of this group")) {
			t.Errorf("Expected error message not found in response: %s", rr.Body.String())
		}
	})
}
