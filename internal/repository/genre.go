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
func (r *Repo) CreateGenre(ctx context.Context, gnr models.Genre) (models.Genre, error) {
	query := "INSERT INTO genres (name) VALUES (?)"

	res, err := r.DB.ExecContext(ctx, query, gnr.Name)
	if err != nil {
		return models.Genre{}, fmt.Errorf("executing insertion to genres table: %w", err)
	}

	gnr.ID, err = res.LastInsertId()
	if err != nil {
		return models.Genre{}, fmt.Errorf("getting inserted ID while adding genre: %w", err)
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
	query := `DELETE FROM genres WHERE id = ?;`

	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting genre: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while deleting genre: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
