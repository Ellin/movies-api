package models

import "database/sql"

// Repo represents the application's access to the database. Used in the repository layer.
type Repo struct {
	DB *sql.DB
}

type Genre struct {
	ID   int
	Name string
}

type Movie struct {
	ID          int
	Title       string
	ReleaseYear int
	Duration    int
}

type Actor struct {
	ID        int
	Name      string
	BirthDate string
}
