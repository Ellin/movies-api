package repository

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

// Repo represents the application's access to the database. Used in the repository layer.
type Repo struct {
	DB *sql.DB
}

// isForeignKeyError checks if the error is a foreign key constraint error.
func isForeignKeyError(err error) bool {
	var sqliteErr sqlite3.Error

	// Check if err is a foreign key constraint error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
		return true
	}

	return false
}
