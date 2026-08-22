package service

import (
	"context"
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type MovieService struct {
	repo     *repository.Repo
	validate *validator.Validate
}

type MovieSubmission struct {
	Title       string  `json:"title" validate:"required"`
	ReleaseYear int     `json:"release_year" validate:"required"`
	Duration    int     `json:"duration" validate:"required,gte=1,lte=100000"`
	GenreIDs    []int64 `json:"genre_ids" validate:"dive,gte=1"`
	ActorIDs    []int64 `json:"actor_ids" validate:"dive,gte=1"`
}

// MoviePatch uses pointers so users can do partial updates for movie data
// Nil pointer values can be used to distinguish data not provided from zero/empty values
type MoviePatch struct {
	Title       *string  `json:"title"`
	ReleaseYear *int     `json:"release_year"`
	Duration    *int     `json:"duration" validate:"omitempty,gte=1,lte=100000"`
	GenreIDs    *[]int64 `json:"genre_ids" validate:"omitempty,dive,gte=1"`
	ActorIDs    *[]int64 `json:"actor_ids" validate:"omitempty,dive,gte=1"`
}

func NewMovieService(r *repository.Repo, v *validator.Validate) *MovieService {
	return &MovieService{repo: r, validate: v}
}

// For movies, your service should allow adding new movies with their title, release year, duration, associated genre, and actors.
func (ms *MovieService) AddMovie(ctx context.Context, sub MovieSubmission) (models.Movie, error) {
	// Struct level validation
	if err := ms.validate.Struct(sub); err != nil {
		return models.Movie{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	// Business level validations
	if err := validateReleaseYear(sub.ReleaseYear); err != nil {
		return models.Movie{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	newMovie := models.Movie{
		Title:       sub.Title,
		ReleaseYear: sub.ReleaseYear,
		Duration:    sub.Duration,
		GenreIDs:    sub.GenreIDs,
		ActorIDs:    sub.ActorIDs,
	}

	// add movie into movies table
	movie, err := ms.repo.AddMovie(ctx, newMovie)

	// returned movie includes a newly generated id
	return movie, err
}

// You'll need functions to retrieve all movies, fetch a specific movie by ID, and filter movies by genre or release year.

func (ms *MovieService) GetMovie(ctx context.Context, id int64) (models.MovieDetail, error) {
	if id < 1 {
		return models.MovieDetail{}, errors.New("id must be positive")
	}

	movie, err := ms.repo.GetMovie(ctx, id)

	return movie, err
}

func (ms *MovieService) GetAllMovies(ctx context.Context) ([]models.MovieDetail, error) {
	return ms.repo.GetAllMovies(ctx)
}

func (ms *MovieService) PatchMovie(ctx context.Context, id int64, patch MoviePatch) (models.Movie, error) {
	// Struct level validation
	if err := ms.validate.Struct(patch); err != nil {
		return models.Movie{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	// First get existing movie
	movieDetail, err := ms.GetMovie(ctx, id)
	if err != nil {
		return models.Movie{}, err
	}

	movie := stripMovieDetails(movieDetail)

	// Update non-nil values of user input
	if patch.Title != nil {
		movie.Title = *patch.Title
	}

	if patch.ReleaseYear != nil {
		//  Validate release year
		if err := validateReleaseYear(*patch.ReleaseYear); err != nil {
			return models.Movie{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
		}
		movie.ReleaseYear = *patch.ReleaseYear
	}

	if patch.Duration != nil {
		movie.Duration = *patch.Duration
	}

	if patch.GenreIDs != nil {
		movie.GenreIDs = *patch.GenreIDs
	}

	if patch.ActorIDs != nil {
		movie.ActorIDs = *patch.ActorIDs
	}

	// Update database with updated movie
	movie, err = ms.repo.PatchMovie(ctx, movie)
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

// stripMovieDetails converts a MovieDetail object to Movie (removing genre and actor names)
func stripMovieDetails(md models.MovieDetail) models.Movie {
	movie := models.Movie{
		ID:          md.ID,
		Title:       md.Title,
		ReleaseYear: md.ReleaseYear,
		Duration:    md.Duration,
	}

	// Strip genre and actor details (leaving only genre and actor IDs)
	for _, genre := range md.Genres {
		movie.GenreIDs = append(movie.GenreIDs, genre.ID)
	}

	for _, actor := range md.Actors {
		movie.ActorIDs = append(movie.ActorIDs, actor.ID)
	}

	return movie
}

func (ms *MovieService) DeleteMovie(ctx context.Context, id int64) error {
	err := ms.repo.DeleteMovie(ctx, id)
	fmt.Println(err)
	return err
}

func (ms *MovieService) GetMoviesByYear(ctx context.Context, year int) ([]models.MovieDetail, error) {
	if err := validateReleaseYear(year); err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	movies, err := ms.repo.GetMoviesByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

// validateReleaseYear checks that the movie's release year is between 1888 and the current year
func validateReleaseYear(year int) error {
	earliestMovie := 1888
	currentYear := time.Now().Year()

	if year < earliestMovie {
		return fmt.Errorf("release year cannot be earlier than %v", earliestMovie)
	}
	if year > currentYear {
		return fmt.Errorf("release year cannot be in the future")
	}

	return nil
}
