package models

import (
	"testing"
	"time"
)

func TestGroupModel(t *testing.T) {
	group := Group{
		ID:          1,
		Title:       "Test Group",
		Description: "This is a test group.",
		CreatorID:   101,
		CreatedAt:   time.Now(),
	}

	if group.ID != 1 {
		t.Errorf("Expected ID 1, got %d", group.ID)
	}
	if group.Title != "Test Group" {
		t.Errorf("Expected Title 'Test Group', got %s", group.Title)
	}
	if group.Description != "This is a test group." {
		t.Errorf("Expected Description 'This is a test group.', got %s", group.Description)
	}
	if group.CreatorID != 101 {
		t.Errorf("Expected CreatorID 101, got %d", group.CreatorID)
	}
}
