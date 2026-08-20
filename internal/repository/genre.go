package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/models"
)

//CRUD for genres

// CREATE
func (r *Repo) AddGenre(ctx context.Context, gnr models.Genre) (models.Genre, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Genre{}, fmt.Errorf("opening transaction for genre addition: %w", err)
	}

	defer tx.Rollback()

	//add genre into genres table
	query := "INSERT INTO genres (name) VALUES (?)"

	res, err := r.DB.ExecContext(ctx, query, gnr.Name)
	if err != nil {
		return models.Genre{}, fmt.Errorf("executing insertion to genres table: %w", err)
	}

	//get the id of freshly added genre
	gnr.ID, err = res.LastInsertId()
	if err != nil {
		return models.Genre{}, fmt.Errorf("getting inserted ID while adding genre: %w", err)
	}

	//add connections to JOIN table movies_genres

	for _, movieID := range gnr.MovieIDs {
		query := "INSERT INTO movies_genres (genre_id, movie_id) VALUES (?, ?)"

		_, err := tx.ExecContext(ctx, query, gnr.ID, movieID)
		if err != nil {
			return models.Genre{}, fmt.Errorf("making connection in movies_genres table: %w", err)
		}

	}
	//commit transaction

	err = tx.Commit()
	if err != nil {
		return models.Genre{}, fmt.Errorf("committing transaction for adding genre: %w", err)
	}

	return gnr, nil
}

// READ
func (r *Repo) GetGenre(ctx context.Context, id int64) (models.Genre, error) {
	query := "SELECT * FROM genres WHERE id = ?"

	row := r.DB.QueryRowContext(ctx, query, id)
	genre := models.Genre{}

	err := row.Scan(&genre.ID, &genre.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Genre{}, ErrNotFound
		}
		return models.Genre{}, fmt.Errorf("scanning the data from genre row to struct: %w", err)
	}

	return genre, nil

}

// READ 1.2
func (r *Repo) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	query := "SELECT id, name FROM genres ORDER BY name"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting genres from genre table: %w", err)
	}

	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
		genre := models.Genre{}

		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating genre rows while getting all genres: %w", err)
	}

	return genres, nil
}

// UPDATE
func (r *Repo) PatchGenre(ctx context.Context, g models.Genre) (models.Genre, error) {
	query := `UPDATE genres
	SET name = ?
	WHERE id = ?;`
	result, err := r.DB.ExecContext(ctx, query, g.Name, g.ID)
	if err != nil {
		return models.Genre{}, fmt.Errorf("updating genre: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return models.Genre{}, fmt.Errorf("getting affected rows while updating genre: %w", err)
	}

	if rows == 0 {
		return models.Genre{}, ErrNotFound
	}

	return g, nil
}

//DELETE

func (r *Repo) DeleteGenre(ctx context.Context, id int64) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("deliting from genre table: %w", err)
	}

	defer tx.Rollback()

	//delete from genres table
	query := `DELETE FROM genres WHERE id = ?;`

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting genre: %w", err)
	}

	//getting the affected row to validate that we actually deleted the row
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while deleting genre: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	//remove all the connections with the genre that is being removed
	connectionQuery := `DELETE FROM genres_movies WHERE genre_id = ?`
	_, err = tx.ExecContext(ctx, connectionQuery, id)
	if err != nil {
		return fmt.Errorf("deleting connection from genres_movies: %w", err)
	}

	//commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}
	return nil
}
