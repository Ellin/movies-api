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

// PostMovie decodes and creates a new movie from the request body,
// responding with the created movie on success.
func (app *App) PostMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub := service.MovieSubmission{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid request body", errs.ErrInvalidUserInput))
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

// GetAllMovies returns a paginated list of movies based on query parameters (e.g. genre, year).
// Filters are parsed via parseFilters before being passed to the service.
func (app *App) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()

	// Check filters from query
	filter, parseErr := parseFilters(query)
	if parseErr != nil {
		errs.WriteError(w, parseErr)
		return
	}

	movies, totalCount, err := app.MovieService.GetAllMovies(ctx, filter)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	response := PaginatedResponse[models.MovieDetail]{
		Data:       movies,
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalCount: totalCount,
		TotalPages: calcTotalPages(totalCount, filter.Pagination.PageSize),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetMovie returns the movie matching the {id} path parameter.
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

// GetMovieSearch returns a paginated list of movies whose title matches the "title" query parameter.
// The search is a case-insensitive, partial match search.
func (app *App) GetMovieSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	titleSearch := r.URL.Query().Get("title")

	pagData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	movies, totalCount, err := app.MovieService.GetMovieSearch(ctx, titleSearch, pagData)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	response := PaginatedResponse[models.MovieDetail]{
		Data:       movies,
		Page:       pagData.Page,
		PageSize:   pagData.PageSize,
		TotalCount: totalCount,
		TotalPages: calcTotalPages(totalCount, pagData.PageSize),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// PatchMovie decodes the request body and updates the movie matching the {id} path parameter
// with the provided fields. Responds with the updated movie on success.
func (app *App) PatchMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	patch := service.MoviePatch{}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid request body", errs.ErrInvalidUserInput))
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

// DeleteMovie removes the movie matching the {id} path parameter.
func (app *App) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	var force bool = r.URL.Query().Get("force") == "true"

	if err := app.MovieService.DeleteMovie(ctx, id, force); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMovieActors returns a list of actors featured in the movie matching the {id} path parameter.
func (app *App) GetMovieActors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	movieID, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	pageData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	actors, totalCount, err := app.MovieService.GetActorsByMovie(ctx, movieID, pageData)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	response := PaginatedResponse[models.ActorSummary]{
		Data:       actors,
		Page:       pageData.Page,
		PageSize:   pageData.PageSize,
		TotalCount: totalCount,
		TotalPages: calcTotalPages(totalCount, pageData.PageSize),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// parseFilters parses all filter parameters from query values into models.MovieFilter
func parseFilters(query url.Values) (models.MovieFilter, error) {
	var filter models.MovieFilter

	// Parse "year"
	if queryYear := query.Get("year"); queryYear != "" {
		year, err := strconv.Atoi(queryYear)
		if err != nil {
			return models.MovieFilter{}, fmt.Errorf("%w: year must be an integer", errs.ErrInvalidUserInput)
		}

		filter.ReleaseYear = &year
	}

	// Parse "genre"
	if queryGenre := query.Get("genre"); queryGenre != "" {
		genre, err := parseID(queryGenre)
		if err != nil {
			return models.MovieFilter{}, fmt.Errorf("%w: invalid genre: %w", errs.ErrInvalidUserInput, err)
		}

		filter.Genre = &genre
	}

	// Parse "actor"
	if queryActor := query.Get("actor"); queryActor != "" {
		actor, err := parseID(queryActor)
		if err != nil {
			return models.MovieFilter{}, fmt.Errorf("%w: invalid actor: %w", errs.ErrInvalidUserInput, err)
		}

		filter.Actor = &actor
	}

	// Parse pagination input
	pagination, err := pagination.Parse(query)
	if err != nil {
		return models.MovieFilter{}, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}
	filter.Pagination = pagination

	return filter, nil
}

// parseID parses an ID string into an int64
func parseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64) // int64 equivalent of Atoi
	if err != nil || id < 1 {
		return 0, fmt.Errorf("id must be positive integer")
	}

	return id, nil
}
