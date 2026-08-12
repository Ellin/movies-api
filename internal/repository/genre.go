package repository

import (
	"database/sql"
	"errors"
	"movies-api/internal/models"
)

//CRUD for genres

// CREATE
func (r *Repo) CreateGenre(gnr models.Genre) (int64, error) {
	query := "INSERT INTO genres (name) VALUES (?)"

	res, err := r.DB.Exec(query, gnr.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// READ
func (r *Repo) GetGenreByID(id int64) (*models.Genre, error) {
	query := "SELECT * FROM genres WHERE id = ?"

	row := r.DB.QueryRow(query, id)
	genre := &models.Genre{}

	err := row.Scan(&genre.ID, &genre.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, err
	}

	return genre, nil

}

// READ 1.2
func (r *Repo) GetAllGenres() ([]*models.Genre, error) {
	query := "SELECT * FROM genres ORDER BY name"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		genre := &models.Genre{}

		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)

	}

	return genres, nil
}
