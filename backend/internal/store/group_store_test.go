package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	createTableSQL := `CREATE TABLE groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		creator_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		privacy TEXT NOT NULL DEFAULT 'public'
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestCreateGroup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewGroupStore(db)

	group := &models.Group{
		CreatorID:   1,
		Title:       "Test Group",
		Description: "This is a test group.",
	}

	createdGroup, err := store.CreateGroup(group)
	if err != nil {
		t.Fatalf("CreateGroup failed: %v", err)
	}

	if createdGroup.ID == 0 {
		t.Error("Expected created group to have an ID")
	}

	var title string
	err = db.QueryRow("SELECT title FROM groups WHERE id = ?", createdGroup.ID).Scan(&title)
	if err != nil {
		t.Fatalf("Failed to query created group: %v", err)
	}

	if title != "Test Group" {
		t.Errorf("Expected group title to be 'Test Group', got '%s'", title)
	}
}
