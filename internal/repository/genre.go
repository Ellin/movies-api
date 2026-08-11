package repository

import (
	"errors"
	"movies-api/internal/models"
)

var ErrorNotFound = errors.New("record not found")

func (r *Repo) CreateGenre(gnr models.Genre) (int64, error) {
	query := "INSERT INTO genres (name) VALUES (?)"

	res, err := r.DB.Exec(query, gnr.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repo) GetGenreByID(id int64) (*models.Genre, error) {
	query := "SELECT * FROM genres WHERE id = ?"

	row := r.DB.QueryRow(query, id)
	g := &models.Genre{}

	err := row.Scan(&g.ID, &g.Name)
	if err != nil {
		return nil, err
	}

	return g, nil

}
