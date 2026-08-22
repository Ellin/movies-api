package service

import (
	"context"
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

type ActorService struct {
	repo     *repository.Repo
	validate *validator.Validate
}

type ActorSubmission struct {
	Name      string  `json:"name" validate:"required"`
	BirthDate string  `json:"birth_date" validate:"required"`
	MovieIDs  []int64 `json:"movie_ids" validate:"dive,gte=1"`
}

// ActorPatch uses pointers so users can do partial updates for movie data
// Nil pointer values can be used to distinguish data not provided from zero/empty values
type ActorPatch struct {
	Name      *string  `json:"name"`
	BirthDate *string  `json:"birth_date"`
	MovieIDs  *[]int64 `json:"movie_ids" validate:"omitempty,dive,gte=1"`
}

var earliestActorBirthDate = time.Date(
	1914, 11, 8,
	0, 0, 0, 0,
	time.UTC,
)

func NewActorService(r *repository.Repo, v *validator.Validate) *ActorService {
	return &ActorService{
		repo:     r,
		validate: v,
	}
}

// For actors, your service should allow adding new actors with their name, birth date and associated movies.
func (as *ActorService) AddActor(ctx context.Context, sub ActorSubmission) (models.Actor, error) {

	if err := as.validate.Struct(sub); err != nil {
		return models.Actor{}, fmt.Errorf(
			"%w: %w",
			errs.ErrInvalidUserInput,
			err,
		)
	}

	if err := validateActorName(sub.Name); err != nil {
		return models.Actor{}, fmt.Errorf(
			"%w: %w",
			errs.ErrInvalidUserInput,
			err,
		)
	}

	if err := validateBirthDate(sub.BirthDate); err != nil {
		return models.Actor{}, fmt.Errorf(
			"%w: %w",
			errs.ErrInvalidUserInput,
			err,
		)
	}

	newActor := models.Actor{
		Name:      strings.TrimSpace(sub.Name),
		BirthDate: sub.BirthDate,
		MovieIDs:  sub.MovieIDs,
	}

	return as.repo.AddActor(ctx, newActor)
}

func validateBirthDate(birthDate string) error {
	date, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return errors.New("birth date must have format YYYY-MM-DD")
	}

	if date.Before(earliestActorBirthDate) {
		return fmt.Errorf(
			"birth date cannot be earlier than %s",
			earliestActorBirthDate,
		)
	}

	if date.After(time.Now()) {
		return errors.New("birth date cannot be in the future")
	}

	return nil
}

func validateActorName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("name cannot be empty")
	}

	if utf8.RuneCountInString(name) > 100 {
		return errors.New("name is too long")
	}

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || r == '-' || r == '\'' {
			continue
		}

		return errors.New("name contains invalid characters")
	}

	return nil
}

// You'll need functions to retrieve all actors, fetch a specific actor by ID, and filter actors by movie or birth date.

func (as *ActorService) GetActor(ctx context.Context, id int64) (models.Actor, error) {

	if id < 1 {
		return models.Actor{}, fmt.Errorf(
			"%w: id must be positive",
			errs.ErrInvalidUserInput,
		)
	}

	return as.repo.GetActor(ctx, id)
}

func (as *ActorService) GetAllActors(ctx context.Context) ([]models.Actor, error) {
	return as.repo.GetAllActors(ctx)
}

func (as *ActorService) PatchActor(ctx context.Context, id int64, patch ActorPatch) (models.Actor, error) {
	// Struct level validation
	if err := as.validate.Struct(patch); err != nil {
		return models.Actor{}, fmt.Errorf(
			"%w: %w",
			errs.ErrInvalidUserInput,
			err,
		)
	}

	// Get existing actor
	actor, err := as.GetActor(ctx, id)
	if err != nil {
		return models.Actor{}, err
	}

	// Update name if provided
	if patch.Name != nil {
		if err := validateActorName(*patch.Name); err != nil {
			return models.Actor{}, fmt.Errorf(
				"%w: %w",
				errs.ErrInvalidUserInput,
				err,
			)
		}

		actor.Name = strings.TrimSpace(*patch.Name)
	}

	// Update birth date if provided
	if patch.BirthDate != nil {
		if err := validateBirthDate(*patch.BirthDate); err != nil {
			return models.Actor{}, fmt.Errorf(
				"%w: %w",
				errs.ErrInvalidUserInput,
				err,
			)
		}

		actor.BirthDate = *patch.BirthDate
	}

	// Replace movie relationships if provided
	if patch.MovieIDs != nil {
		actor.MovieIDs = *patch.MovieIDs
	}

	return as.repo.PatchActor(ctx, actor)
}

func (as *ActorService) DeleteActor(ctx context.Context, id int64) error {

	if id < 1 {
		return fmt.Errorf(
			"%w: id must be positive",
			errs.ErrInvalidUserInput,
		)
	}

	return as.repo.DeleteActor(ctx, id)
}
