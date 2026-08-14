package main

import (
	"encoding/json"
	"errors"
	"movies-api/internal/repository"
	"net/http"
	"strconv"
)

func (app *application) postMovie(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {

}

// get by ID
func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64) // int64 equivalent of Atoi
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	movie, err := app.movieService.GetMovie(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
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
