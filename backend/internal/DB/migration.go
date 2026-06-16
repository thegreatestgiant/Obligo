package DB

import (
	"database/sql"
	"embed"
	"io/fs"
	// ...
)

//go:embed schema/*.sql
var migrationFiles embed.FS

func RunMigrations(db *sql.DB) error {
	content, err := fs.ReadFile(migrationFiles, "schema/000_official_look.sql")
	if err != nil {
		return err
	}

	// Execute the SQL against the database
	_, err = db.Exec(string(content))
	return err
}
