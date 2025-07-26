package websocket

import (
	"testing"
)

func TestManagerRegisterUnregister(t *testing.T) {
	manager := NewManager(nil, nil, nil)

	client := NewClient(123, "test-client-123", nil)

	// Test registration
	manager.Register(client)
	if !manager.IsOnline(123) {
		t.Error("Client should be online after registration")
	}

	// Test unregistration
	manager.Unregister(123)
	if manager.IsOnline(123) {
		t.Error("Client should be offline after unregistration")
	}
}

func TestManagerBroadcast(t *testing.T) {
	manager := NewManager(nil, nil, nil)

	// Create mock clients
	client1 := NewClient(1, "test-client-1", nil)
	client2 := NewClient(2, "test-client-2", nil)

	manager.Register(client1)
	manager.Register(client2)

	// Test online status
	if !manager.IsOnline(1) || !manager.IsOnline(2) {
		t.Error("Clients should be online after registration")
	}

	// Test getting online user IDs
	ids := manager.OnlineUserIDs()
	if len(ids) != 2 {
		t.Errorf("Expected 2 online users, got %d", len(ids))
	}
}
