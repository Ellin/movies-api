package service

import (
	"context"
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/pagination"
	"movies-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type MovieService struct {
	repo     *repository.Repo
	validate *validator.Validate
}

type MovieSubmission struct {
	Title       string  `json:"title" validate:"required,max=255"`
	ReleaseYear int     `json:"release_year" validate:"required"`
	Duration    int     `json:"duration" validate:"required,gte=1,lte=100000"`
	GenreIDs    []int64 `json:"genre_ids" validate:"dive,gte=1"`
	ActorIDs    []int64 `json:"actor_ids" validate:"dive,gte=1"`
}

// MoviePatch uses pointers so users can do partial updates for movie data
// Nil pointer values can be used to distinguish data not provided from zero/empty values
type MoviePatch struct {
	Title       *string  `json:"title" validate:"omitempty,min=1,max=255"`
	ReleaseYear *int     `json:"release_year"`
	Duration    *int     `json:"duration" validate:"omitempty,gte=1,lte=100000"`
	GenreIDs    *[]int64 `json:"genre_ids" validate:"omitempty,dive,gte=1"`
	ActorIDs    *[]int64 `json:"actor_ids" validate:"omitempty,dive,gte=1"`
}

func NewMovieService(r *repository.Repo, v *validator.Validate) *MovieService {
	return &MovieService{repo: r, validate: v}
}

// For movies, your service should allow adding new movies with their title, release year, duration, associated genre, and actors.
func (ms *MovieService) AddMovie(ctx context.Context, sub MovieSubmission) (models.MovieDetail, error) {
	// Struct level validation
	if err := ms.validate.Struct(sub); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return models.MovieDetail{}, handleValidationError(validationErrors)
		}
		return models.MovieDetail{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	// Business level validations
	if err := validateReleaseYear(sub.ReleaseYear); err != nil {
		return models.MovieDetail{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}
	for _, actorID := range sub.ActorIDs {
		actor, err := ms.repo.GetActor(ctx, actorID)
		if err != nil {
			return models.MovieDetail{}, err
		}

		if err := validateActorMovieDates(actor.BirthDate, sub.ReleaseYear); err != nil {
			return models.MovieDetail{}, fmt.Errorf(
				"%w: %w",
				errs.ErrInvalidUserInput,
				err,
			)
		}
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

func (ms *MovieService) GetMovie(ctx context.Context, id int64) (models.MovieDetail, error) {
	if id < 1 {
		return models.MovieDetail{}, errors.New("id must be positive")
	}

	movie, err := ms.repo.GetMovie(ctx, id)

	return movie, err
}

// GetAllMovies returns movies matching filter's criteria, restricted to the pagination parameters.
// Also returns the total number of matching movies across all pages.
func (ms *MovieService) GetAllMovies(ctx context.Context, filter models.MovieFilter) ([]models.MovieDetail, int, error) {
	if err := ms.validateFilter(filter); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	return ms.repo.GetAllMovies(ctx, filter)
}

// validateFilter is a helper that validates all query parameters
func (ms *MovieService) validateFilter(filter models.MovieFilter) error {
	if err := filter.Pagination.Validate(); err != nil {
		return err
	}

	if filter.ReleaseYear != nil {
		if err := validateReleaseYear(*filter.ReleaseYear); err != nil {
			return err
		}
	}

	if filter.Actor != nil {
		if *filter.Actor < 1 {
			return fmt.Errorf("actor ID must be positive")
		}
	}

	if filter.Genre != nil {
		if *filter.Genre < 1 {
			return fmt.Errorf("genre ID must be positive")
		}
	}

	return nil
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

func (ms *MovieService) PatchMovie(ctx context.Context, id int64, patch MoviePatch) (models.MovieDetail, error) {
	// Struct level validation
	if err := ms.validate.Struct(patch); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return models.MovieDetail{}, handleValidationError(validationErrors)
		}
		return models.MovieDetail{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	// First get existing movie
	movieDetail, err := ms.GetMovie(ctx, id)
	if err != nil {
		return models.MovieDetail{}, err
	}

	movie := stripMovieDetails(movieDetail)

	// Update non-nil values of user input
	if patch.Title != nil {
		movie.Title = *patch.Title
	}

	if patch.ReleaseYear != nil {
		//  Validate release year
		if err := validateReleaseYear(*patch.ReleaseYear); err != nil {
			return models.MovieDetail{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
		}
		//  Validate that all participating actors birth dates
		for _, actorID := range movie.ActorIDs {
			actor, err := ms.repo.GetActor(ctx, actorID)
			if err != nil {
				return models.MovieDetail{}, err
			}

			if err := validateActorMovieDates(actor.BirthDate, *patch.ReleaseYear); err != nil {
				return models.MovieDetail{}, fmt.Errorf(
					"%w: %w",
					errs.ErrInvalidUserInput,
					err,
				)
			}
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
		for _, actorID := range *patch.ActorIDs {
			actor, err := ms.repo.GetActor(ctx, actorID)
			if err != nil {
				return models.MovieDetail{}, err
			}
			if err := validateActorMovieDates(actor.BirthDate, movie.ReleaseYear); err != nil {
				return models.MovieDetail{}, fmt.Errorf(
					"%w: %w",
					errs.ErrInvalidUserInput,
					err,
				)
			}
		}
		movie.ActorIDs = *patch.ActorIDs
	}

	// Update database with updated movie
	return ms.repo.PatchMovie(ctx, movie)
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

func (ms *MovieService) DeleteMovie(ctx context.Context, id int64, force bool) error {
	return ms.repo.DeleteMovie(ctx, id, force)
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

func (ms *MovieService) GetActorsByMovie(ctx context.Context, movieID int64, pageData pagination.Pagination) ([]models.ActorSummary, int, error) {
	if movieID < 1 {
		return nil, 0, fmt.Errorf("%w: movie ID must be positive", errs.ErrInvalidUserInput)
	}
	return ms.repo.GetActorsByMoviePaginated(ctx, movieID, pageData)
}
