// For interacting with database (CRUD operations)
package repository

import (
	"database/sql"
	"errors"
	"movies-api/internal/apperrors"
	"movies-api/internal/models"
)

// AddMovie inserts a new movie into the movies table (CREATE)
func (r *Repo) AddMovie(m models.Movie) (int64, error) {
	query := `INSERT INTO movies (title, releaseYear, duration)
	VALUES (?, ?, ?);`

	result, err := r.DB.Exec(query, m.Title, m.ReleaseYear, m.Duration)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetMovie gets movie data from the movies table (READ)
func (r *Repo) GetMovie(id int64) (models.Movie, error) {
	query := `SELECT id, title, releaseYear, duration
	FROM movies
	WHERE id = ?;`

	var m models.Movie
	err := r.DB.QueryRow(query, id).Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration) // fill movie struct with data from found row
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Movie{}, apperrors.ErrNoRecord
		} else {
			return models.Movie{}, err
		}

	}

	return m, nil
}

// UPDATE: Update movie

// DELETE: Delete movie
