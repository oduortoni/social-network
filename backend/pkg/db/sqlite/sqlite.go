package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Migration runs all .sql files in the migrations/sqlite directory to set up the database schema
func Migration() (*sql.DB, error) {
	db, err := CheckDB()
	if err != nil {
		return nil, err
	}

	// 1. Ensure the schema_migrations table exists
	if err := ensureSchemaVersionTable(db); err != nil {
		return nil, fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	// 2. Get all applied migrations
	appliedMigrations, err := GetAppliedMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// 3. Find and run new migrations
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	migrationDir := filepath.Join(wd, "pkg/db/migrations/sqlite")
	log.Printf("Using migration directory: %s", migrationDir)

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration dir: %w", err)
	}

	// Sort files to ensure they are applied in order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	migrated := false
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		// Check if migration was already applied
		if _, ok := appliedMigrations[file.Name()]; ok {
			continue
		}

		log.Printf("Executing migration: %s", file.Name())

		if err := ApplyMigrationInTx(db, migrationDir, file.Name()); err != nil {
			return nil, fmt.Errorf("failed to apply migration %s: %w", file.Name(), err)
		}

		log.Printf("Migrated: %s", file.Name())
		migrated = true
	}

	if !migrated {
		log.Printf("Database is up to date. No new migrations to apply.")
	}
	return db, nil
}
