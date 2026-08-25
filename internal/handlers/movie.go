package handlers

import (
	"encoding/json"
	"fmt"
	"movies-api/internal/errs"
	"movies-api/internal/models"
	"movies-api/internal/pagination"
	"movies-api/internal/service"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

// parseFilters parses all filter parameters from query values into models.MovieFilter
func parseFilters(query url.Values) (models.MovieFilter, error) {
	var filter models.MovieFilter

	// Parse "year"
	if queryYear := query.Get("year"); queryYear != "" {
		year, err := strconv.Atoi(queryYear)
		if err != nil {
			return models.MovieFilter{}, errs.ErrInvalidUserInput // invalid input
		}

		filter.ReleaseYear = &year
	}

	// Parse "genre"
	if queryGenre := query.Get("genre"); queryGenre != "" {
		genre, err := parseID(queryGenre)
		if err != nil {
			return models.MovieFilter{}, errs.ErrInvalidUserInput // invalid input
		}

		filter.Genre = &genre
	}

	// Parse pagination input
	pagination, err := pagination.Parse(query)
	if err != nil {
		return models.MovieFilter{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}
	filter.Pagination = pagination

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
	if filter.ReleaseYear != nil {
		movies, err = app.MovieService.GetMoviesByYear(ctx, *filter.ReleaseYear)
	} else if filter.Genre != nil {
		movies, err = app.MovieService.GetMoviesByGenre(ctx, *filter.Genre)
	} else {
		// Get all movies (no filters)
		movies, err = app.MovieService.GetAllMovies(ctx, filter)
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
