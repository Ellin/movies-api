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
	Name *string `json:"name" validate:"omitempty,min=1,max=255"`
}

// initialize genre service
func NewGenreService(r *repository.Repo, v *validator.Validate) *GenreService {
	return &GenreService{repo: r, validate: v}
}

// validate the data and call the repo lvl
func (gs *GenreService) AddGenre(ctx context.Context, sub GenreSubmission) (models.Genre, error) {

	if err := gs.validate.Struct(sub); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return models.Genre{}, handleValidationError(validationErrors)
		}
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

func (gs *GenreService) GetAllGenres(ctx context.Context, pag pagination.Pagination) ([]models.GenreSummary, int, error) {
	if err := pag.Validate(); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}
	return gs.repo.GetAllGenres(ctx, pag)
}

func (gs *GenreService) GetGenre(ctx context.Context, id int64) (models.Genre, error) {
	if id < 1 {
		return models.Genre{}, errors.New("id must be positive")
	}
	return gs.repo.GetGenre(ctx, id)
}

func (ms *MovieService) GetMovieSearch(ctx context.Context, title string, pag pagination.Pagination) ([]models.MovieDetail, int, error) {
	if err := pag.Validate(); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}
	movies, totalcount, err := ms.repo.GetMovieSearch(ctx, title, pag)
	if err != nil {
		return nil, 0, err
	}

	return movies, totalcount, nil
}

func (gs *GenreService) PatchGenre(ctx context.Context, id int64, patch GenrePatch) (models.GenreSummary, error) {
	g := models.GenreSummary{ID: id}

	if err := gs.validate.Struct(patch); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return models.GenreSummary{}, handleValidationError(validationErrors)
		}
		return models.GenreSummary{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
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
