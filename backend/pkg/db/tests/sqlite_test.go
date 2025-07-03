package tests

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigration(t *testing.T) {
	testDBPath := "test_social_network.db"
	os.Setenv("SQLITE_DB_PATH", testDBPath)
	defer os.Remove(testDBPath)
	defer os.Unsetenv("SQLITE_DB_PATH")
}
