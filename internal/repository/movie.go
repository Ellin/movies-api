// For interacting with database (CRUD operations)
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/apperrors"
	"movies-api/internal/models"
)

// AddMovie inserts a new movie into the movies table (CREATE)
func (r *Repo) AddMovie(m models.Movie) (int64, error) {
	query := `INSERT INTO movies (title, releaseYear, duration)
	VALUES (?, ?, ?);`

	result, err := r.DB.Exec(query, m.Title, m.ReleaseYear, m.Duration)
	if err != nil {
		return 0, fmt.Errorf("adding movie: %w", err)
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
			return models.Movie{}, fmt.Errorf("getting movie: %w", err)
		}
	}

	return m, nil
}

// UpdateMovie updates the movie with matching ID and returns the number of rows affected (UPDATE)
func (r *Repo) UpdateMovie(m models.Movie) (int64, error) {
	query := `UPDATE movies
	SET title = ?, releaseYear = ?, duration = ?
	WHERE id = ?;`

	result, err := r.DB.Exec(query, m.Title, m.ReleaseYear, m.Duration, m.ID)
	if err != nil {
		return 0, fmt.Errorf("updating movie: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("getting affected rows while updating movie: %w", err)
	}

	if rows == 0 {
		return 0, apperrors.ErrNoRecord
	}

	return rows, nil
}

// DeleteMovie deletes the movie with match ID and returns the number of rows affected (DELETE)
func (r *Repo) DeleteMovie(id int64) (int64, error) {
	query := `DELETE FROM movies
	where id = ?;`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("deleting movie: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("getting affected rows while deleting movie: %w", err)
	}

	if rows == 0 {
		return 0, apperrors.ErrNoRecord
	}

	return rows, nil
}
