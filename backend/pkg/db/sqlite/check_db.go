package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// checks if the database file exists, and creates it if not
func CheckDB() (*sql.DB, error) {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "social_network.db"
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create db file: %w", err)
		}
		file.Close()
		log.Printf("Created new database file: %s", dbPath)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	return db, nil
}
