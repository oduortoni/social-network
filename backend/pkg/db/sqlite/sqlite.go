package sqlite

import (
	"database/sql"

	migration "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migration.NewWithDatabaseInstance(
		"file://backend/pkg/db/migrations/sqlite",
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}
	return m.Up()
}
