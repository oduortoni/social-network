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
	appliedMigrations, err := getAppliedMigrations(db)
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

		path := filepath.Join(migrationDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		log.Printf("Executing migration: %s", file.Name())

		// Run migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction for migration %s: %w", file.Name(), err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("migration failed for %s: %w", file.Name(), err)
		}

		// Record the migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", file.Name()); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to record migration %s: %w", file.Name(), err)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction for migration %s: %w", file.Name(), err)
		}

		log.Printf("Migrated: %s", file.Name())
		migrated = true
	}

	if !migrated {
		log.Printf("Database is up to date. No new migrations to apply.")
	}
	return db, nil
}

func ensureSchemaVersionTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}
