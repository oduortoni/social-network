package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)



// Migration runs all .sql files in the migrations/sqlite directory to set up the database schema
func Migration() (*sql.DB, error) {
	db, err := CheckDB()
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	migrationDir := wd + "/pkg/db/migrations/sqlite"
	log.Printf("Using migration directory: %s", migrationDir)
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration dir: %w", err)
	}

	migrated := false
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}
		path := fmt.Sprintf("%s/%s", migrationDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read migration file %s: %v", file.Name(), err)
			continue
		}
		log.Printf("Executing migration: %s", file.Name())
		if _, err := db.Exec(string(content)); err != nil {
			log.Printf("Migration failed for %s: %v", file.Name(), err)
			continue
		}
		log.Printf("Migrated: %s", file.Name())
		migrated = true
	}
	if !migrated {
		log.Printf("No migration files found or executed.")
	}
	return db, nil
}