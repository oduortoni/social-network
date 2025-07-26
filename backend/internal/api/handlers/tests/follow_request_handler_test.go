package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// MockFollowRequestService is a mock implementation of the FollowRequestServiceInterface for testing.
type MockFollowRequestService struct {
	AcceptedFollowConnectionFunc func(followConnectionID int64) error
	RejectedFollowConnectionFunc func(followConnectionID int64) error
}

func (s *MockFollowRequestService) AcceptedFollowConnection(followConnectionID int64) error {
	if s.AcceptedFollowConnectionFunc != nil {
		return s.AcceptedFollowConnectionFunc(followConnectionID)
	}
	return nil
}

func (s *MockFollowRequestService) RejectedFollowConnection(followConnectionID int64) error {
	if s.RejectedFollowConnectionFunc != nil {
		return s.RejectedFollowConnectionFunc(followConnectionID)
	}
	return nil
}

func TestFollowRequestRespond_AcceptSuccess(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{
		AcceptedFollowConnectionFunc: func(followConnectionID int64) error {
			if followConnectionID == 123 {
				return nil
			}
			return sql.ErrNoRows
		},
	}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "accepted"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/123/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "123")

	// Add user ID to context (simulating authenticated user)
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Successfully accepted follow request"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_RejectSuccess(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{
		RejectedFollowConnectionFunc: func(followConnectionID int64) error {
			if followConnectionID == 123 {
				return nil
			}
			return sql.ErrNoRows
		},
	}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "rejected"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/123/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "123")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Successfully rejected follow request"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_RequestNotFound(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{
		AcceptedFollowConnectionFunc: func(followConnectionID int64) error {
			return sql.ErrNoRows // Request not found
		},
	}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "accepted"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/999/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "999")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Follow request not found"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_InvalidRequestID(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "accepted"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/invalid/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "invalid")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Invalid request ID"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_InvalidStatus(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "invalid"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/123/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "123")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Invalid status. Must be 'accepted' or 'rejected'"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_InvalidJSON(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/follow-request/123/request", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "123")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Invalid JSON format"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_Unauthorized(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "accepted"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/123/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "123")

	// Don't add user ID to context (simulating unauthenticated user)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "User not found in context"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}

func TestFollowRequestRespond_ServiceError(t *testing.T) {
	mockFollowRequestService := &MockFollowRequestService{
		AcceptedFollowConnectionFunc: func(followConnectionID int64) error {
			return sql.ErrConnDone // Simulate database connection error
		},
	}
	followRequestHandler := handlers.NewFollowRequestHandler(mockFollowRequestService)

	requestStatus := models.FollowRequestResponseStatus{Status: "accepted"}
	body, _ := json.Marshal(requestStatus)
	req := httptest.NewRequest("POST", "/follow-request/123/request", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("requestId", "123")

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	followRequestHandler.FollowRequestRespond(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Failed to accept follow request"
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}
