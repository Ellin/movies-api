// For interacting with database (CRUD operations)
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/models"
)

// AddMovie inserts a new movie into the movies table (CREATE)
func (r *Repo) AddMovie(ctx context.Context, m models.Movie) (models.Movie, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Movie{}, fmt.Errorf("beginning transaction for adding movie: %w", err)
	}
	defer tx.Rollback()

	// Insert into movies table
	query := `INSERT INTO movies (title, releaseYear, duration) VALUES (?, ?, ?);`
	result, err := tx.ExecContext(ctx, query, m.Title, m.ReleaseYear, m.Duration)
	if err != nil {
		return models.Movie{}, fmt.Errorf("adding movie: %w", err)
	}

	// Get newly inserted movie ID
	m.ID, err = result.LastInsertId()
	if err != nil {
		return models.Movie{}, fmt.Errorf("getting inserted ID while adding movie: %w", err)
	}

	// Insert into genres_movies table
	for _, genreID := range m.GenreIDs {
		query := `INSERT INTO genres_movies (genre_id, movie_id) VALUES (?, ?);`
		_, err = tx.ExecContext(ctx, query, genreID, m.ID)
		if err != nil {
			return models.Movie{}, fmt.Errorf("linking genre %d to movie: %w", genreID, err)
		}
	}

	// Insert into movies_actors table
	for _, actorID := range m.ActorIDs {
		query := `INSERT INTO movies_actors (movie_id, actor_id) VALUES (?, ?);`
		_, err = tx.ExecContext(ctx, query, m.ID, actorID)
		if err != nil {
			return models.Movie{}, fmt.Errorf("linking actor %d to movie: %w", actorID, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.Movie{}, fmt.Errorf("commiting transaction for adding movie: %w", err)
	}

	return m, nil
}

// GetMovie gets movie data from the movies table (READ)
func (r *Repo) GetMovie(ctx context.Context, id int64) (models.Movie, error) {
	query := `SELECT id, title, releaseYear, duration
	FROM movies
	WHERE id = ?;`

	var m models.Movie
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration) // fill movie struct with data from found row
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Movie{}, ErrNotFound
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
		return 0, ErrNotFound
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
		return 0, ErrNotFound
	}

	return rows, nil
}
