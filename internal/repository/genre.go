package repository

import (
	"errors"
	"movies-api/internal/models"
)

var ErrorNotFound = errors.New("record not found")

func (r *Repo) CreateGenre(gnr models.Genre) (int64, error) {
	query := "INSERT INTO genre (name) VALUES (?)"

	res, err := r.DB.Exec(query, gnr.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
