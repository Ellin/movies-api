package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/models"
)

// AddActor inserts a new actor into the actors table (CREATE)
func (r *Repo) AddActor(ctx context.Context, a models.Actor) (models.Actor, error) {
	// Create new transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Actor{}, fmt.Errorf("beginning transaction for adding actor: %w", err)
	}
	defer tx.Rollback()

	// Insert into actors table
	query := `INSERT INTO actors (name, birth_date) VALUES (?, ?);`
	result, err := tx.ExecContext(ctx, query, a.Name, a.BirthDate)
	if err != nil {
		return models.Actor{}, fmt.Errorf("adding actor: %w", err)
	}

	// Get newly inserted actor ID
	a.ID, err = result.LastInsertId()
	if err != nil {
		return models.Actor{}, fmt.Errorf("getting inserted ID while adding actor: %w", err)
	}

	// Insert into movies_actors table
	for _, movieID := range a.MovieIDs {
		query := `INSERT INTO movies_actors (movie_id, actor_id) VALUES (?, ?);`
		_, err = tx.ExecContext(ctx, query, movieID, a.ID)
		if err != nil {
			return models.Actor{}, fmt.Errorf("linking actor %d to movie: %w", movieID, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.Actor{}, fmt.Errorf("commiting transaction for adding actor: %w", err)
	}

	return a, nil
}

// GetActor gets actor data from the actors table (READ)
func (r *Repo) GetActor(ctx context.Context, id int64) (models.ActorDetail, error) {
	// Get data from actors table
	query := `SELECT id, name, birth_date
	FROM actors WHERE id = ?;`

	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return models.ActorDetail{}, err
	}
	defer rows.Close()

	var a models.ActorDetail
	err = r.DB.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Name, &a.BirthDate) // fill actor struct with data from found row
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ActorDetail{}, ErrNotFound
		} else {
			return models.ActorDetail{}, fmt.Errorf("getting actor: %w", err)
		}
	}

	a.Movies, err = r.getMoviesByActor(ctx, id)
	if err != nil {
		return models.ActorDetail{}, err
	}

	return a, nil
}

// getMoviesByActor is a helper that retrieves all movies (with id and title) associated with a given actor ID
func (r *Repo) getMoviesByActor(ctx context.Context, actorID int64) ([]models.MovieDetail, error) {
	// Get actor data associated with the given movie ID
	query := `SELECT ma.movie_id, m.title
	FROM movies_actors ma JOIN movies m ON ma.movie_id = m.id
	WHERE ma.actor_id = ?;`

	rows, err := r.DB.QueryContext(ctx, query, actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieDetail
	for rows.Next() {
		var movie models.MovieDetail
		err = rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Genres, &movie.Actors)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}
		movies = append(movies, movie)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return movies, nil
}

func (r *Repo) GetAllActors(ctx context.Context) ([]models.ActorDetail, error) {
	// make movie-actors map with actor id as key
	actorMoviesMap, err := r.buildActorMoviesMap(ctx)
	if err != nil {
		return nil, err
	}

	// get all actors from actors table
	query := `SELECT id, name, birth_date FROM actors ORDER BY id ASC;`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all actors: %w", err)
	}
	defer rows.Close()

	var allActors []models.ActorDetail

	// Get actors row by row
	for rows.Next() {
		var a models.ActorDetail
		err = rows.Scan(&a.ID, &a.Name, &a.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("scanning actor row while getting all actors: %w", err)
		}

		a.Movies = actorMoviesMap[a.ID]

		allActors = append(allActors, a)
	}

	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating actor rows while getting all actors: %w", err)
	}

	return allActors, nil
}

// buildActorMoviesMap is a helper that creates a map where the key is the actor ID and value is all associated movies
func (r *Repo) buildActorMoviesMap(ctx context.Context) (map[int64][]models.MovieDetail, error) {
	// join movies_actors with movies
	query := `SELECT actor_id, movie_id, m.title AS movie_title, m.release_year AS movie_release_year, movie.duration AS movie_duration
	FROM movies_actors ma JOIN movies m on m.id = ma.movie_id
	ORDER BY actor_id ASC;`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting actor-movies rows: %w", err)
	}
	defer rows.Close()

	// make movie-actors map with movie id as key
	actorMoviesMap := make(map[int64][]models.MovieDetail)
	for rows.Next() {
		var actorID int64
		var movie models.MovieDetail

		err = rows.Scan(&actorID, &movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Duration)
		if err != nil {
			return nil, fmt.Errorf("scanning movie actor row while getting all movies: %w", err)
		}

		actorMoviesMap[actorID] = append(actorMoviesMap[actorID], movie)
	}
	// Check if rows.Next() loop stopped due to error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie actor rows while getting all movies: %w", err)
	}

	return actorMoviesMap, nil
}

// PatchActor updates the actors table along with relationships in movies_actors table with updated actor information
func (r *Repo) PatchActor(ctx context.Context, a models.Actor) (models.Actor, error) {
	// Create new transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Actor{}, fmt.Errorf("beginning transaction for patching actor: %w", err)
	}
	defer tx.Rollback()

	// Update actors table
	if err := r.updateActor(ctx, tx, a); err != nil {
		return models.Actor{}, err
	}

	// Update movies_actors table
	if err := r.updateMoviesActors(ctx, tx, a); err != nil {
		return models.Actor{}, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.Actor{}, fmt.Errorf("commiting transaction for patching actor: %w", err)
	}

	return a, nil
}

// updateActor is a helper that updates the actors table with matching ID
func (r *Repo) updateActor(ctx context.Context, tx *sql.Tx, a models.Actor) error {
	// Update actors table
	query := `UPDATE actors
	SET name = ?, birth_date = ?
	WHERE id = ?;`

	result, err := tx.ExecContext(ctx, query, a.Name, a.BirthDate, a.ID)
	if err != nil {
		return fmt.Errorf("updating actor: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while updating actor: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// updateMoviesActors is a helper that updates the movies_actors table
func (r *Repo) updateMoviesActors(ctx context.Context, tx *sql.Tx, a models.Actor) error {
	// Delete all existing relations between the given actor ID and movies in movies_actors table
	query := `DELETE FROM movies_actors WHERE actor_id = ?;`
	_, err := tx.ExecContext(ctx, query, a.ID)
	if err != nil {
		return err
	}

	// Add new relationships
	for _, movieID := range a.MovieIDs {
		query := `INSERT INTO movies_actors(movie_id, actor_id) VALUES(?, ?)
		ON CONFLICT(movie_id, actor_id) DO NOTHING;`

		_, err := tx.ExecContext(ctx, query, movieID, a.ID)
		if err != nil {
			return fmt.Errorf("linking actor %d to movie: %w", movieID, err)
		}
	}

	return nil
}

// DeleteActor deletes the actor by ID (DELETE)
func (r *Repo) DeleteActor(ctx context.Context, id int64) error {
	query := `DELETE FROM actors WHERE id = ?;`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting actor: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while deleting actor: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
