package service

import (
	"context"
	"errors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type MovieService struct {
	repo *repository.Repo
}

func NewMovieService(r *repository.Repo) *MovieService {
	return &MovieService{repo: r}
}

type MovieSubmission struct {
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	Duration    int    `json:"duration"`
	GenreIDs    []int  `json:"genre_ids"`
	ActorIDs    []int  `json:"actor_ids"`
}

// For movies, your service should allow adding new movies with their title, release year, duration, associated genre, and actors.
func (ms *MovieService) AddMovie(ctx context.Context, sub MovieSubmission) (models.Movie, error) {
	//  Validate release year

	// Validate duration

	// validate genreIDs

	// Validate ActorIDs

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

func (ms *MovieService) GetMovie(ctx context.Context, id int64) (models.Movie, error) {
	if id < 1 {
		return models.Movie{}, errors.New("id must be positive")
	}

	movie, err := ms.repo.GetMovie(ctx, id)

	return movie, err
}

// Don't forget to implement a way to get all actors in a specific movie.

// Updating a movie's details (including its title, release year, duration, genre, and actors) and removing a movie should also be supported.
