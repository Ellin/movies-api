package service

import (
	"context"
	"errors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

// service layer for genres
type GenreService struct {
	repo *repository.Repo
}

// package prepared for submission
type GenreSubmission struct {
	Name string `json: name`
}

// initialize genre service
func GenreNewService(r *repository.Repo) *GenreService {
	return &GenreService{repo: r}
}

// validate the data and call the repo lvl
func (gs *GenreService) AddGenre(ctx context.Context, sub GenreSubmission) (models.Genre, error) {
	//pack data in
	NewGenre := models.Genre{
		Name: sub.Name,
	}

	//call repo lvl
	g, err := gs.repo.CreateGenre(ctx, NewGenre)

	return g, err
}

func (gs *GenreService) GetAllGenres(ctx context.Context) ([]*models.Genre, error) {
	return gs.repo.GetAllGenres(ctx)
}

func (gs *GenreService) GetGenre(ctx context.Context, id int64) (*models.Genre, error) {
	if id < 1 {
		return &models.Genre{}, errors.New("id must be positive")
	}
	return gs.repo.GetGenre(ctx, id)
}
