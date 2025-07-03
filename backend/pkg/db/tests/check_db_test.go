package tests

import (
	"os"
	"testing"

	"github.com/tajjjjr/social-network/backend/pkg/db/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func TestCheckDB_CreatesFileIfMissing(t *testing.T) {
	testDBPath := "test_social_network.db"

	if _, err := os.Stat(testDBPath); err == nil {
		os.Remove(testDBPath)
	}

	os.Setenv("SQLITE_DB_PATH", testDBPath)
	defer os.Unsetenv("SQLITE_DB_PATH")

	db, err := sqlite.CheckDB()
	if err != nil {
		t.Fatalf("CheckDB returned an error: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(testDBPath); os.IsNotExist(err) {
		t.Errorf("Expected database file to be created, but it doesn't exist")
	}

	if err := os.Remove(testDBPath); err != nil {
		t.Logf("Warning: Failed to delete test db file: %v", err)
	}
}
