package main

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

func (app *application) postMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub := service.MovieSubmission{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	movie, err := app.movieService.AddMovie(ctx, sub)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before add movie finished")
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)

			// ! TO DO - UPDATE ERROR HANDLING: distinguish between invalid/bad requests and internal server errors. Do not expose internal server error messages.
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	movies, err := app.movieService.GetAllMovies(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before get movie finished")
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// get by ID
func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64) // int64 equivalent of Atoi
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	movie, err := app.movieService.GetMovie(ctx, id)
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
	json.NewEncoder(w).Encode(movie)
}

// update by ID
func (app *application) updateMovie(w http.ResponseWriter, r *http.Request) {

}

// delete by ID
func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {

}
