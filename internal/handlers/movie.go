package handlers

import (
	"encoding/json"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MovieFilter struct {
	releaseYear *int
	genre       *int64
	actor       *int64
}

func (app *App) PostMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub := service.MovieSubmission{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sub.Title = strings.TrimSpace(sub.Title)

	movie, err := app.MovieService.AddMovie(ctx, sub)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// parseFilters parses filters from the query into MovieFilter
func parseFilters(query url.Values) (MovieFilter, error) {
	var filter MovieFilter

	// Parse "year"
	if queryYear := query.Get("year"); queryYear != "" {
		year, err := strconv.Atoi(queryYear)
		if err != nil {
			return MovieFilter{}, errs.ErrInvalidUserInput // invalid input
		}

		filter.releaseYear = &year
	}

	// Parse "genre"
	if queryGenre := query.Get("genre"); queryGenre != "" {
		genre, err := parseID(queryGenre)
		if err != nil {
			return MovieFilter{}, errs.ErrInvalidUserInput // invalid input
		}

		filter.genre = &genre
	}

	if queryActor := query.Get("actor"); queryActor != "" {
		actor, err := parseID(queryActor)
		if err != nil {
			return MovieFilter{}, errs.ErrInvalidUserInput
		}

		filter.actor = &actor
	}

	return filter, nil
}

func parseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64) // int64 equivalent of Atoi
	if err != nil || id < 1 {
		return 0, fmt.Errorf("id must be positive integer")
	}

	return id, nil
}

func (app *App) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()

	// Check filters from query
	filter, parseErr := parseFilters(query)
	if parseErr != nil {
		errs.WriteError(w, parseErr)
		return
	}

	var movies []models.MovieDetail
	var err error

	// Filter movies by release year
	if filter.releaseYear != nil {
		movies, err = app.MovieService.GetMoviesByYear(ctx, *filter.releaseYear)
	} else if filter.genre != nil {
		movies, err = app.MovieService.GetMoviesByGenre(ctx, *filter.genre)
	} else if filter.actor != nil {
		movies, err = app.MovieService.GetMoviesByActor(ctx, *filter.actor)
	} else {
		// Get all movies (no filters)
		movies, err = app.MovieService.GetAllMovies(ctx)
	}

	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

// get by ID
func (app *App) GetMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	movie, err := app.MovieService.GetMovie(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

// update by ID
func (app *App) PatchMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	patch := service.MoviePatch{}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if patch.Title != nil {
		*patch.Title = strings.TrimSpace(*patch.Title)
	}

	movie, err := app.MovieService.PatchMovie(ctx, id, patch)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

// delete by ID
func (app *App) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	if err := app.MovieService.DeleteMovie(ctx, id); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
