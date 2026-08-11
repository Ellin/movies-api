package repository

import "database/sql"

// Repo represents the application's access to the database. Used in the repository layer.
type Repo struct {
	DB *sql.DB
}
