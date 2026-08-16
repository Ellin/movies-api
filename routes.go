package main

import (
	"movies-api/internal/handlers"
	"net/http"
)

func NewRouter(app *handlers.App) *http.ServeMux {
	mux := http.NewServeMux()

	// movie routes
	mux.HandleFunc("POST /api/movie", app.PostMovie)
	mux.HandleFunc("GET /api/movie", app.GetAllMovies)
	mux.HandleFunc("GET /api/movie/{id}", app.GetMovie)
	mux.HandleFunc("PATCH /api/movie/{id}", app.PatchMovie)
	mux.HandleFunc("DELETE /api/movie/{id}", app.DeleteMovie)

	return mux
}
