package repository

import (
	"database/sql"
	"errors"
)

// Repo represents the application's access to the database. Used in the repository layer.
type Repo struct {
	DB *sql.DB
}

// Unified error for non-existing records
var ErrorNotFound = errors.New("record not found")
