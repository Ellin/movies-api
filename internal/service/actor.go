package service

import (
	"context"
	"errors"
	"fmt"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type ActorService struct {
	repo *repository.Repo
}

type ActorSubmission struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	BirthDate string  `json:"birth_date"`
	MovieIDs  []int64 `json:"movie_ids"`
}

// MoviePatch uses pointers so users can do partial updates for movie data
// Nil pointer values can be used to distinguish data not provided from zero/empty values
type ActorPatch struct {
	Name      *string  `json:"name"`
	BirthDate *string  `json:"birth_date"`
	MovieIDs  *[]int64 `json:"movie_ids"`
}

func NewActorService(r *repository.Repo) *ActorService {
	return &ActorService{repo: r}
}

// For movies, your service should allow adding new movies with their title, release year, duration, associated genre, and actors.
func (as *ActorService) AddActor(ctx context.Context, sub ActorSubmission) (models.Actor, error) {
	//  Validate Name

	// Validate Birth Date

	// Validate MovieIDs

	newActor := models.Actor{
		Name:      sub.Name,
		BirthDate: sub.BirthDate,
		MovieIDs:  sub.MovieIDs,
	}

	// add actor into actors table
	actor, err := as.repo.AddActor(ctx, newActor)

	// returned actor includes a newly generated id
	return actor, err
}

// You'll need functions to retrieve all actors, fetch a specific actor by ID, and filter actors by movie or birth date.

func (as *ActorService) GetActor(ctx context.Context, id int64) (models.ActorSummary, error) {
	if id < 1 {
		return models.ActorSummary{}, errors.New("id must be positive")
	}

	actor, err := as.repo.GetActor(ctx, id)

	return actor, err
}

func (as *ActorService) GetAllActors(ctx context.Context) ([]models.ActorSummary, error) {
	return as.repo.GetAllActors(ctx)
}

func (as *ActorService) PatchActor(ctx context.Context, id int64, patch ActorPatch) (models.Actor, error) {
	// First get existing actor
	actorSummary, err := as.GetActor(ctx, id)
	if err != nil {
		return models.Actor{}, err
	}

	actor := stripActorSummary(actorSummary)

	// Update non-nil values of user input
	if patch.Name != nil {
		actor.Name = *patch.Name
	}

	if patch.BirthDate != nil {
		actor.BirthDate = *patch.BirthDate
	}

	if patch.MovieIDs != nil {
		actor.MovieIDs = *patch.MovieIDs
	}

	// Update database with updated actor
	actor, err = as.repo.PatchActor(ctx, actor)
	if err != nil {
		return models.Actor{}, err
	}

	return actor, nil
}

// stripActorSummary converts a ActorSummary object to Actor (removing movie details)
func stripActorSummary(as models.ActorSummary) models.Actor {
	actor := models.Actor{
		ID:        as.ID,
		Name:      as.Name,
		BirthDate: as.BirthDate,
	}

	// Strip movie details (leaving only movie IDs)
	for _, movie := range as.Movies {
		actor.MovieIDs = append(actor.MovieIDs, movie.ID)
	}

	return actor
}

func (as *ActorService) DeleteActor(ctx context.Context, id int64) error {
	err := as.repo.DeleteActor(ctx, id)
	fmt.Println(err)
	return err
}
