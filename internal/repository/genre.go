package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/pagination"
)

//CRUD for genres

// CREATE
func (r *Repo) AddGenre(ctx context.Context, gnr models.Genre) (models.Genre, error) {
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
			return models.Genre{}, errs.ErrNotFound
		}
		return models.Genre{}, fmt.Errorf("scanning the data from genre row to struct: %w", err)
	}

	genre.MovieIDs, err = r.buildMovieIDslice(ctx, id)
	if err != nil {
		return models.Genre{}, fmt.Errorf("getting movie IDs for genre: %w", err)
	}

	return genre, nil

}

// READ 1.2
func (r *Repo) GetAllGenres(ctx context.Context, pag pagination.Pagination) ([]models.GenreSummary, error) {
	query := "SELECT id, name FROM genres ORDER BY name ASC LIMIT ? OFFSET ?"
	rows, err := r.DB.QueryContext(ctx, query, pag.Limit(), pag.Offset())
	if err != nil {
		return nil, fmt.Errorf("getting genres from genre table: %w", err)
	}

	defer rows.Close()

	var genres []models.GenreSummary

	for rows.Next() {
		genre := models.GenreSummary{}

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

func (r *Repo) buildMovieIDslice(ctx context.Context, gID int64) ([]int64, error) {
	query := `SELECT movie_id FROM genres_movies WHERE genre_id = ?`

	rows, err := r.DB.QueryContext(ctx, query, gID)
	if err != nil {
		return nil, fmt.Errorf("getting rows from genres_movies: %w", err)
	}
	var movieIDs []int64

	for rows.Next() {
		var movie int64

		err = rows.Scan(&movie)
		if err != nil {
			return nil, fmt.Errorf("scanning row in genres_movies: %w", err)
		}
		movieIDs = append(movieIDs, movie)
	}

	return movieIDs, nil
}

// UPDATE
func (r *Repo) PatchGenre(ctx context.Context, g models.GenreSummary) (models.GenreSummary, error) {
	query := `UPDATE genres
	SET name = ?
	WHERE id = ?;`
	result, err := r.DB.ExecContext(ctx, query, g.Name, g.ID)
	if err != nil {
		return models.GenreSummary{}, fmt.Errorf("updating genre: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return models.GenreSummary{}, fmt.Errorf("getting affected rows while updating genre: %w", err)
	}

	if rows == 0 {
		return models.GenreSummary{}, errs.ErrNotFound
	}

	return g, nil
}

//DELETE

func (r *Repo) DeleteGenre(ctx context.Context, id int64, force bool) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("deliting from genre table: %w", err)
	}

	defer tx.Rollback()
	//check if there're connections
	checkConnectionQuery := `SELECT movie_id FROM genres_movies WHERE genre_id = ?`
	rows, err := tx.QueryContext(ctx, checkConnectionQuery, id)
	if err != nil {
		return fmt.Errorf("getting rows from genres_movies: %w", err)
	}
	if rows != nil && !force {
		return errs.ErrForce
	}

	//delete from genres table
	query := `DELETE FROM genres WHERE id = ?;`

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting genre: %w", err)
	}

	//getting the affected row to validate that we actually deleted the row
	rowNum, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while deleting genre: %w", err)
	}

	if rowNum == 0 {
		return errs.ErrNotFound
	}

	//commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}
	return nil
}

// // getMoviesByGenre is a helper that retrieves all movies (with id and title) associated with a given actor ID
func (r *Repo) getMoviesByGenre(ctx context.Context, genreID int64) ([]models.MovieDetail, error) {
	// Get genre data associated with the given movie ID
	query := `SELECT ma.movie_id, m.title
	FROM genres_movies ma JOIN movies m ON ma.movie_id = m.id
	WHERE ma.genre_id = ?;`

	rows, err := r.DB.QueryContext(ctx, query, genreID)
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
