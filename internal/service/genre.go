package service

import (
	"context"
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/pagination"
	"movies-api/internal/repository"

	"github.com/go-playground/validator/v10"
)

// service layer for genres
type GenreService struct {
	repo     *repository.Repo
	validate *validator.Validate
}

// package prepared for submission
type GenreSubmission struct {
	Name string `json:"name" validate:"required"`
}

type GenrePatch struct {
	Name *string `json:"name"`
}

// initialize genre service
func NewGenreService(r *repository.Repo, v *validator.Validate) *GenreService {
	return &GenreService{repo: r, validate: v}
}

// validate the data and call the repo lvl
func (gs *GenreService) AddGenre(ctx context.Context, sub GenreSubmission) (models.Genre, error) {

	if err := gs.validate.Struct(sub); err != nil {
		return models.Genre{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	//pack data in
	NewGenre := models.Genre{
		Name: sub.Name,
	}

	//call repo lvl
	g, err := gs.repo.AddGenre(ctx, NewGenre)

	return g, err
}

func (gs *GenreService) GetAllGenres(ctx context.Context, pag pagination.Pagination) ([]models.GenreSummary, error) {
	if err := pag.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}
	return gs.repo.GetAllGenres(ctx, pag)
}

func (gs *GenreService) GetGenre(ctx context.Context, id int64) (models.Genre, error) {
	if id < 1 {
		return models.Genre{}, errors.New("id must be positive")
	}
	return gs.repo.GetGenre(ctx, id)
}

func (gs *GenreService) PatchGenre(ctx context.Context, id int64, patch GenrePatch) (models.GenreSummary, error) {
	g := models.GenreSummary{ID: id}
	if patch.Name != nil {
		g.Name = *patch.Name
	}

	g, err := gs.repo.PatchGenre(ctx, g)
	if err != nil {
		return models.GenreSummary{}, err
	}

	return g, nil
}

func (gs *MovieService) DeleteGenre(ctx context.Context, id int64, force bool) error {
	err := gs.repo.DeleteGenre(ctx, id, force)
	return err
}
