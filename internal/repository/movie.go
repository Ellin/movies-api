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

	// ! TO DO: ALSO GET GENRE AND ACTORS INFO

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Movie{}, ErrNotFound
		} else {
			return models.Movie{}, fmt.Errorf("getting movie: %w", err)
		}
	}

	return m, nil
}

func (r *Repo) GetAllMovies(ctx context.Context) ([]models.MovieDetail, error) {
	// get all movies from movies table
	query := `SELECT id, title, releaseYear, duration FROM movies ORDER BY id ASC;`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all movies: %w", err)
	}
	defer rows.Close()

	// ---- GENRE-MOVIES

	// join genres_movies data with genre names
	query = `SELECT movie_id, genre_id, g.name AS genre_name
	FROM genres_movies gm JOIN genres g ON g.id = gm.genre_id
	ORDER BY movie_id ASC;`
	mgrows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting movie_genre rows: %w", err)
	}
	defer mgrows.Close()

	// make movie-genres map with movie id as key
	movieGenresMap := make(map[int64][]models.Genre)
	for mgrows.Next() {
		var movieID int64
		var genre models.Genre

		err = mgrows.Scan(&movieID, &genre.ID, &genre.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning movie genre row while getting all movies: %w", err)
		}

		movieGenresMap[movieID] = append(movieGenresMap[movieID], genre)
	}

	// Check if mgrows.Next() loop stopped due to error
	if err := mgrows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie genre rows while getting all movies: %w", err)
	}

	// ---- MOVIES-ACTORS
	// join movies_actors with actors
	query = `SELECT movie_id, actor_id, a.name AS actor_name
	FROM movies_actors ma JOIN actors a on a.id = ma.actor_id
	ORDER BY movie_id ASC;`
	marows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting movie-actors rows: %w", err)
	}
	defer marows.Close()

	// make movie-actors map with movie id as key
	movieActorsMap := make(map[int64][]models.ActorSummary)
	for marows.Next() {
		var movieID int64
		var actor models.ActorSummary

		err = marows.Scan(&movieID, &actor.ID, &actor.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning movie actor row while getting all movies: %w", err)
		}

		movieActorsMap[movieID] = append(movieActorsMap[movieID], actor)
	}
	// Check if magrows.Next() loop stopped due to error
	if err := marows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie actor rows while getting all movies: %w", err)
	}

	var allMovies []models.MovieDetail

	// Get movies row by row
	for rows.Next() {
		var m models.MovieDetail
		err = rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("scanning movie row while getting all movies: %w", err)
		}

		// add genres info
		m.Genres = movieGenresMap[m.ID]

		// add actors info
		m.Actors = movieActorsMap[m.ID]

		allMovies = append(allMovies, m)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie rows while getting all movies: %w", err)
	}

	return allMovies, nil
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
