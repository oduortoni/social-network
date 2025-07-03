package tests

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigration(t *testing.T) {
	testDBPath := "test_social_network.db"
	os.Setenv("SQLITE_DB_PATH", testDBPath)
	defer os.Remove(testDBPath)
	defer os.Unsetenv("SQLITE_DB_PATH")


	testDir := "/pkg/db/migrations/sqlite"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll("pkg")

	upMigration := `
	CREATE TABLE IF NOT EXISTS test_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`
	upFile := filepath.Join(testDir, "001_create_test_users.up.sql")
	if err := os.WriteFile(upFile, []byte(upMigration), 0644); err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}
}
