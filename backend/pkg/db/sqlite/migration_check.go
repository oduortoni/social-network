package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func ensureSchemaVersionTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
	return err
}

func GetAppliedMigrations(db *sql.DB) (map[string]bool, error) {
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

	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return applied, nil
}

func ApplyMigrationInTx(db *sql.DB, migrationDir, fileName string) error {
	path := filepath.Join(migrationDir, fileName)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// Apply the migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	// Record the migration
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", fileName); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit the transaction
	return tx.Commit()
}
