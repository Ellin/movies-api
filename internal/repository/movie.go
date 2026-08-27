// For interacting with database (CRUD operations)
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"strings"
)

// AddMovie inserts a new movie into the movies table (CREATE)
func (r *Repo) AddMovie(ctx context.Context, m models.Movie) (models.MovieDetail, error) {
	// Create new transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.MovieDetail{}, fmt.Errorf("beginning transaction for adding movie: %w", err)
	}
	defer tx.Rollback()

	// Insert into movies table
	query := `INSERT INTO movies (title, releaseYear, duration) VALUES (?, ?, ?);`
	result, err := tx.ExecContext(ctx, query, m.Title, m.ReleaseYear, m.Duration)
	if err != nil {
		return models.MovieDetail{}, fmt.Errorf("adding movie: %w", err)
	}

	// Get newly inserted movie ID
	m.ID, err = result.LastInsertId()
	if err != nil {
		return models.MovieDetail{}, fmt.Errorf("getting inserted ID while adding movie: %w", err)
	}

	// Insert into genres_movies table
	for _, genreID := range m.GenreIDs {
		query := `INSERT INTO genres_movies (genre_id, movie_id) VALUES (?, ?)
		ON CONFLICT (genre_id, movie_id) DO NOTHING;`
		_, err = tx.ExecContext(ctx, query, genreID, m.ID)
		if err != nil {
			if isForeignKeyError(err) {
				return models.MovieDetail{}, fmt.Errorf("%w: referenced genre ID invalid", errs.ErrInvalidUserInput)
			}
			return models.MovieDetail{}, fmt.Errorf("linking genre %d to movie: %w", genreID, err)
		}
	}

	// Insert into movies_actors table
	for _, actorID := range m.ActorIDs {
		query := `INSERT INTO movies_actors (movie_id, actor_id) VALUES (?, ?)
		ON CONFLICT (movie_id, actor_id) DO NOTHING;`
		_, err = tx.ExecContext(ctx, query, m.ID, actorID)
		if err != nil {
			if isForeignKeyError(err) {
				return models.MovieDetail{}, fmt.Errorf("%w: referenced actor ID invalid", errs.ErrInvalidUserInput)
			}
			return models.MovieDetail{}, fmt.Errorf("linking actor %d to movie: %w", actorID, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.MovieDetail{}, fmt.Errorf("commiting transaction for adding movie: %w", err)
	}

	return r.GetMovie(ctx, m.ID)
}

// GetMovie gets movie data from the movies table (READ)
func (r *Repo) GetMovie(ctx context.Context, id int64) (models.MovieDetail, error) {
	query := `SELECT id, title, releaseYear, duration
	FROM movies WHERE id = ?;`

	var m models.MovieDetail
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration) // fill movie struct with data from found row
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.MovieDetail{}, errs.ErrNotFound
		} else {
			return models.MovieDetail{}, fmt.Errorf("getting movie: %w", err)
		}
	}

	m.Genres, err = r.getGenresByMovie(ctx, id)
	if err != nil {
		return models.MovieDetail{}, err
	}

	m.Actors, err = r.GetActorsByMovie(ctx, id)
	if err != nil {
		return models.MovieDetail{}, err
	}

	return m, nil
}

// getGenresByMovie is a helper that retrieves all genres (with id and name) associated with a given movie ID
func (r *Repo) getGenresByMovie(ctx context.Context, movieID int64) ([]models.GenreSummary, error) {
	// Get genre data associated with the given movie ID
	query := `SELECT gm.genre_id, g.name
	FROM genres_movies gm JOIN genres g ON gm.genre_id = g.id
	WHERE gm.movie_id = ?;`

	rows, err := r.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.GenreSummary
	for rows.Next() {
		var g models.GenreSummary
		err = rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}

		genres = append(genres, g)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return genres, nil
}

// getActorsByMovie is a helper that retrieves all actors (with id and name) associated with a given movie ID
func (r *Repo) GetActorsByMovie(ctx context.Context, movieID int64) ([]models.ActorSummary, error) {
	// Get actor data associated with the given movie ID
	query := `SELECT ma.actor_id, a.name
	FROM movies_actors ma JOIN actors a ON ma.actor_id = a.id
	WHERE ma.movie_id = ?;`

	rows, err := r.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.ActorSummary
	for rows.Next() {
		var actor models.ActorSummary
		err = rows.Scan(&actor.ID, &actor.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}
		actors = append(actors, actor)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return actors, nil
}

// GetAllMovies gets all movies according to the provided query parameters. Result is paginated with 0-based page numbering.
func (r *Repo) GetAllMovies(ctx context.Context, filter models.MovieFilter) ([]models.MovieDetail, error) {
	query, args := createQueryFromFilters(filter)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting all movies: %w", err)
	}
	defer rows.Close()

	allMovies, err := r.scanMoviesFromRows(ctx, rows)
	if err != nil {
		return nil, err
	}

	return allMovies, nil
}

// createQueryFromFilters is a helper that constructs a query string for from MovieFilter. Returns the query and arguments for query execution.
func createQueryFromFilters(filter models.MovieFilter) (query string, args []any) {
	// Base query for getting movies from movies table
	query = `SELECT m.id, m.title, m.releaseYear, m.duration FROM movies m`

	var joins []string  // JOIN statement parts
	var wheres []string // WHERE statement parts

	// Check filters and collect JOIN, WHERE, and argument parts for building query statement
	if filter.Genre != nil {
		joins = append(joins, "JOIN genres_movies gm ON m.id = gm.movie_id")
		wheres = append(wheres, "gm.genre_id = ?")
		args = append(args, filter.Genre)
	}

	if filter.Actor != nil {
		joins = append(joins, "JOIN movies_actors ma ON m.id = ma.movie_id")
		wheres = append(wheres, "ma.actor_id = ?")
		args = append(args, filter.Actor)
	}

	if filter.ReleaseYear != nil {
		wheres = append(wheres, "m.releaseYear = ?")
		args = append(args, filter.ReleaseYear)
	}

	// Construct and add the JOIN statements to query
	if len(joins) > 0 {
		query = query + " " + strings.Join(joins, " ")
	}

	// Construct and add the WHERE statement to query
	if len(wheres) > 0 {
		query += " WHERE " + strings.Join(wheres, " AND ")
	}

	// Finalize query with order and pagination
	query += " ORDER BY m.id ASC LIMIT ? OFFSET ?;"
	args = append(args, filter.Pagination.Limit(), filter.Pagination.Offset())

	return query, args
}

// scanMoviesFromRows is a helper that scans rows selected from the movies table and returns []models.MovieDetail
func (r *Repo) scanMoviesFromRows(ctx context.Context, rows *sql.Rows) ([]models.MovieDetail, error) {
	// make movie-genres map with movie id as key
	movieGenresMap, err := r.buildMovieGenresMap(ctx)
	if err != nil {
		return nil, err
	}

	// make movie-actors map with movie id as key
	movieActorsMap, err := r.buildMovieActorsMap(ctx)
	if err != nil {
		return nil, err
	}

	var movies []models.MovieDetail

	// Scan rows and add each movie to movies
	for rows.Next() {
		var m models.MovieDetail
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("scanning movie row: %w", err)
		}

		m.Genres = movieGenresMap[m.ID]
		m.Actors = movieActorsMap[m.ID]

		movies = append(movies, m)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie rows while getting all movies: %w", err)
	}

	return movies, nil
}

// buildMovieGenresMap is a helper that creates a map where the key is the movie ID and value is all associated genres
func (r *Repo) buildMovieGenresMap(ctx context.Context) (map[int64][]models.GenreSummary, error) {
	// join genres_movies data with genre names
	query := `SELECT movie_id, genre_id, g.name AS genre_name
	FROM genres_movies gm JOIN genres g ON g.id = gm.genre_id
	ORDER BY movie_id ASC;`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting movie_genre rows: %w", err)
	}
	defer rows.Close()

	// make movie-genres map with movie id as key
	movieGenresMap := make(map[int64][]models.GenreSummary)
	for rows.Next() {
		var movieID int64
		var genre models.GenreSummary

		err = rows.Scan(&movieID, &genre.ID, &genre.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning movie genre row while getting all movies: %w", err)
		}

		movieGenresMap[movieID] = append(movieGenresMap[movieID], genre)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie genre rows while getting all movies: %w", err)
	}

	return movieGenresMap, nil
}

// buildMovieActorsMap is a helper that creates a map where the key is the movie ID and value is all associated actors
func (r *Repo) buildMovieActorsMap(ctx context.Context) (map[int64][]models.ActorSummary, error) {
	// join movies_actors with actors
	query := `SELECT movie_id, actor_id, a.name AS actor_name
	FROM movies_actors ma JOIN actors a on a.id = ma.actor_id
	ORDER BY movie_id ASC;`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting movie-actors rows: %w", err)
	}
	defer rows.Close()

	// make movie-actors map with movie id as key
	movieActorsMap := make(map[int64][]models.ActorSummary)
	for rows.Next() {
		var movieID int64
		var actor models.ActorSummary

		err = rows.Scan(&movieID, &actor.ID, &actor.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning movie actor row while getting all movies: %w", err)
		}

		movieActorsMap[movieID] = append(movieActorsMap[movieID], actor)
	}
	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie actor rows while getting all movies: %w", err)
	}

	return movieActorsMap, nil
}

// PatchMovie updates the movies table along with relationships in genres_movies and movies_actors tables with updated movie information
func (r *Repo) PatchMovie(ctx context.Context, m models.Movie) (models.MovieDetail, error) {
	// Create new transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.MovieDetail{}, fmt.Errorf("beginning transaction for patching movie: %w", err)
	}
	defer tx.Rollback()

	// Update movies table
	if err := r.updateMovie(ctx, tx, m); err != nil {
		return models.MovieDetail{}, err
	}

	// Update genres_movies table
	if err := r.updateMovieGenres(ctx, tx, m); err != nil {
		return models.MovieDetail{}, err
	}

	// Update movies_actors table
	if err := r.updateMovieActors(ctx, tx, m); err != nil {
		return models.MovieDetail{}, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.MovieDetail{}, fmt.Errorf("commiting transaction for patching movie: %w", err)
	}

	return r.GetMovie(ctx, m.ID)
}

// updateMovie is a helper that updates the movies table with matching ID
func (r *Repo) updateMovie(ctx context.Context, tx *sql.Tx, m models.Movie) error {
	// Update movies table
	query := `UPDATE movies
	SET title = ?, releaseYear = ?, duration = ?
	WHERE id = ?;`

	result, err := tx.ExecContext(ctx, query, m.Title, m.ReleaseYear, m.Duration, m.ID)
	if err != nil {
		return fmt.Errorf("updating movie: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while updating movie: %w", err)
	}

	if rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// updateMovieGenres is a helper that updates the genres_movies table
func (r *Repo) updateMovieGenres(ctx context.Context, tx *sql.Tx, m models.Movie) error {
	// Delete all existing relations between the given movie ID and genres in genres_movies table
	query := `DELETE FROM genres_movies WHERE movie_id = ?;`
	_, err := tx.ExecContext(ctx, query, m.ID)
	if err != nil {
		return err
	}

	// Add new relationships
	for _, genreID := range m.GenreIDs {
		query := `INSERT INTO genres_movies(genre_id, movie_id)	VALUES(?, ?)
		ON CONFLICT(genre_id, movie_id) DO NOTHING;`

		_, err := tx.ExecContext(ctx, query, genreID, m.ID)
		if err != nil {
			if isForeignKeyError(err) {
				return fmt.Errorf("%w: referenced genre ID invalid", errs.ErrInvalidUserInput)
			}
			return fmt.Errorf("linking genre %d to movie: %w", genreID, err)
		}
	}

	return nil
}

// updateMovieActors is a helper that updates the movies_actors table
func (r *Repo) updateMovieActors(ctx context.Context, tx *sql.Tx, m models.Movie) error {
	// Delete all existing relations between the given movie ID and actors in movies_actors table
	query := `DELETE FROM movies_actors WHERE movie_id = ?;`
	_, err := tx.ExecContext(ctx, query, m.ID)
	if err != nil {
		return err
	}

	// Add new relationships
	for _, actorID := range m.ActorIDs {
		query := `INSERT INTO movies_actors(movie_id, actor_id) VALUES(?, ?)
		ON CONFLICT(movie_id, actor_id) DO NOTHING;`

		_, err := tx.ExecContext(ctx, query, m.ID, actorID)
		if err != nil {
			if isForeignKeyError(err) {
				return fmt.Errorf("%w: referenced actor ID invalid", errs.ErrInvalidUserInput)
			}
			return fmt.Errorf("linking actor %d to movie: %w", actorID, err)
		}
	}

	return nil
}

// DeleteMovie deletes the movie by ID. (DELETE)
// Only allows deletion of movies that has associated genres or actors if force is true.
func (r *Repo) DeleteMovie(ctx context.Context, id int64, force bool) error {
	// No force: Check if movie has relationships with genres or actors before deletion.
	if !force {
		var hasRelationships bool

		// Query returns a single row containing 1 (true) or 0 (false)
		query := `SELECT EXISTS (
		SELECT 1 FROM genres_movies WHERE movie_id = ?
		UNION ALL
		SELECT 1 FROM movies_actors WHERE movie_id = ?);`

		if err := r.DB.QueryRowContext(ctx, query, id, id).Scan(&hasRelationships); err != nil {
			return fmt.Errorf("checking movie relationships: %w", err)
		}

		if hasRelationships {
			return errs.ErrForce
		}
	}

	query := `DELETE FROM movies WHERE id = ?;`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting movie: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while deleting movie: %w", err)
	}

	if rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}
