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

	// Parse "actor"
	if queryActor := query.Get("actor"); queryActor != "" {
		actor, err := parseID(queryActor)
		if err != nil {
			return models.MovieFilter{}, errs.ErrInvalidUserInput
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
	var totalCount int
	var err error

	movies, totalCount, err = app.MovieService.GetAllMovies(ctx, filter)
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

	var force bool = r.URL.Query().Get("force") == "true"

	if err := app.MovieService.DeleteMovie(ctx, id, force); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) GetMovieActors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	movieID, err := parseID(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	actors, err := app.MovieService.GetActorsByMovie(ctx, movieID)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
}
