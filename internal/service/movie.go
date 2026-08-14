package service

import (
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
	Title       string
	ReleaseYear int
	Duration    int
	Genres      []int
	Actors      []int
}

// For movies, your service should allow adding new movies with their title, release year, duration, associated genre, and actors.
// func (ms *MovieService) AddMovie(sub MovieSubmission) (models.Movie, error) {

// }

// You'll need functions to retrieve all movies, fetch a specific movie by ID, and filter movies by genre or release year.

func (ms *MovieService) GetMovie(id int64) (models.Movie, error) {
	if id < 1 {
		return models.Movie{}, errors.New("id must be positive")
	}

	movie, err := ms.repo.GetMovie(id)

	return movie, err
}

// Don't forget to implement a way to get all actors in a specific movie.

// Updating a movie's details (including its title, release year, duration, genre, and actors) and removing a movie should also be supported.
