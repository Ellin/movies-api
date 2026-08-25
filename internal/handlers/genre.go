package handlers

import (
	"encoding/json"
	"movies-api/internal/errs"
	"movies-api/internal/pagination"
	"movies-api/internal/service"
	"net/http"
)

// post
func (app *App) PostGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//create DTO (data transfer object) for new genre
	sub := service.GenreSubmission{}

	//parse JSON into DTO
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	//call the service layer
	genre, err := app.GenreService.AddGenre(ctx, sub)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

// get
func (app *App) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, err := pagination.Parse(r.URL.Query())
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	genres, err := app.GenreService.GetAllGenres(ctx, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

// get
func (app *App) GetGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	genres, err := app.GenreService.GetGenre(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (app *App) PatchGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	patch := service.GenrePatch{}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	genre, err := app.GenreService.PatchGenre(ctx, id, patch)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genre)
}

func (app *App) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var force bool = r.URL.Query().Get("force") == "true"

	if err := app.MovieService.DeleteGenre(ctx, id, force); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
