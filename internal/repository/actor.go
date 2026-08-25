package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/errs"
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
		query := `INSERT INTO movies_actors (movie_id, actor_id)
			VALUES (?, ?)
			ON CONFLICT(movie_id, actor_id) DO NOTHING;`

		_, err = tx.ExecContext(ctx, query, movieID, a.ID)
		if err != nil {
			if isForeignKeyError(err) {
				return models.Actor{}, fmt.Errorf(
					"%w: referenced movie ID invalid",
					errs.ErrInvalidUserInput,
				)
			}

			return models.Actor{}, fmt.Errorf(
				"linking movie %d to actor: %w",
				movieID,
				err,
			)
		}
	}
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.Actor{}, fmt.Errorf("commiting transaction for adding actor: %w", err)
	}

	return a, nil
}

// GetActor gets actor data from the actors table (READ)
func (r *Repo) GetActor(ctx context.Context, id int64) (models.Actor, error) {
	query := `SELECT id, name, birth_date FROM actors WHERE id = ?;`

	var actor models.Actor

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&actor.ID, &actor.Name, &actor.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Actor{}, errs.ErrNotFound
		}
		return models.Actor{}, fmt.Errorf("getting actor: %w", err)
	}

	actor.MovieIDs, err = r.getMovieIDsByActor(ctx, id)
	if err != nil {
		return models.Actor{}, err
	}

	return actor, nil
}

// GetAllActors gets all actors data from the actors table (READ)
func (r *Repo) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	actorMoviesMap, err := r.buildActorMovieIDsMap(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, birth_date FROM actors ORDER BY id ASC;`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting all actors: %w", err)
	}
	defer rows.Close()

	var actors []models.Actor

	for rows.Next() {
		var actor models.Actor

		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
			return nil, fmt.Errorf("scanning actor row: %w", err)
		}

		actor.MovieIDs = actorMoviesMap[actor.ID]

		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating actor rows: %w", err)
	}

	return actors, nil
}

func (r *Repo) GetAllActorsByName(ctx context.Context, name string) ([]models.Actor, error) {
	query := `SELECT id, name, birth_date FROM actors WHERE LOWER(name) LIKE LOWER(?) ORDER BY id ASC;`
	rows, err := r.DB.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("getting actors by name: %w", err)
	}
	defer rows.Close()

	var actors []models.Actor
	for rows.Next() {
		var actor models.Actor

		if err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		); err != nil {
			return nil, fmt.Errorf("scanning actor: %w", err)
		}

		actor.MovieIDs, err = r.getMovieIDsByActor(ctx, actor.ID)
		if err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating actors: %w", err)
	}
	return actors, nil
}

// buildActorMoviesMap is a helper that creates a map where the key is the actor ID and value is all associated movies
func (r *Repo) buildActorMovieIDsMap(ctx context.Context) (map[int64][]int64, error) {
	query := `SELECT actor_id, movie_id	FROM movies_actors ORDER BY actor_id, movie_id;`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting actor-movie relationships: %w", err)
	}
	defer rows.Close()

	actorMoviesMap := make(map[int64][]int64)

	for rows.Next() {
		var actorID int64
		var movieID int64

		if err := rows.Scan(&actorID, &movieID); err != nil {
			return nil, fmt.Errorf("scanning actor-movie relationship: %w", err)
		}

		actorMoviesMap[actorID] = append(
			actorMoviesMap[actorID],
			movieID,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating actor-movie relationships: %w", err)
	}

	return actorMoviesMap, nil
}

// getMovieIDsByActor collects movie IDs that are associated with actor ID into slice
func (r *Repo) getMovieIDsByActor(ctx context.Context, actorID int64) ([]int64, error) {
	query := `SELECT movie_id FROM movies_actors WHERE actor_id = ? ORDER BY movie_id;`

	rows, err := r.DB.QueryContext(ctx, query, actorID)
	if err != nil {
		return nil, fmt.Errorf("getting movie IDs for actor: %w", err)
	}
	defer rows.Close()

	var movieIDs []int64

	for rows.Next() {
		var movieID int64

		if err := rows.Scan(&movieID); err != nil {
			return nil, fmt.Errorf("scanning movie ID: %w", err)
		}

		movieIDs = append(movieIDs, movieID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movie IDs: %w", err)
	}

	return movieIDs, nil
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
	query := `UPDATE actors	SET name = ?, birth_date = ? WHERE id = ?;`

	result, err := tx.ExecContext(ctx, query, a.Name, a.BirthDate, a.ID)
	if err != nil {
		return fmt.Errorf("updating actor: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while updating actor: %w", err)
	}

	if rows == 0 {
		return errs.ErrNotFound
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
		query := `INSERT INTO movies_actors(movie_id, actor_id)
			VALUES(?, ?)
			ON CONFLICT(movie_id, actor_id) DO NOTHING;`
		_, err := tx.ExecContext(ctx, query, movieID, a.ID)
		if err != nil {
			if isForeignKeyError(err) {
				return fmt.Errorf(
					"%w: referenced movie ID invalid",
					errs.ErrInvalidUserInput,
				)
			}
			return fmt.Errorf(
				"linking movie %d to actor: %w",
				movieID,
				err,
			)
		}
	}

	return nil
}

// DeleteActor deletes the actor by ID (DELETE)
func (r *Repo) DeleteActor(ctx context.Context, id int64, force bool) error {
	if !force {
		var hasMovies bool
		query := `SELECT EXISTS(SELECT 1 FROM movies_actors WHERE actor_id = ?);`
		if err := r.DB.QueryRowContext(ctx, query, id).Scan(&hasMovies); err != nil {
			return fmt.Errorf("checking actor relationships: %w", err)
		}
		if hasMovies {
			return errs.ErrForce
		}
	}
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
		return errs.ErrNotFound
	}
	return nil
}
