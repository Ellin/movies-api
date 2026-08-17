package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"
	"strconv"
)

// post
func (app *App) PostGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub := service.GenreSubmission{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	genre, err := app.GenreService.AddGenre(ctx, sub)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before add movie finished")
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest) //!!!!!!!!!!!!!
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

// get
func (app *App) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	genres, err := app.GenreService.GetAllGenres(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before get genre finished")
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

// get
func (app *App) GetGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	genres, err := app.GenreService.GetGenre(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before get movie finished")
		} else if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}
