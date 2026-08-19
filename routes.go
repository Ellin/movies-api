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
	//genre routes
	mux.HandleFunc("POST /api/genre", app.PostGenre)
	mux.HandleFunc("GET /api/genre", app.GetAllGenres)
	mux.HandleFunc("GET /api/genre/{id}", app.GetGenre)
	mux.HandleFunc("PATCH /api/movie/{id}", app.PatchGenre)
	mux.HandleFunc("DELETE /api/movie/{id}", app.DeleteGenre)

	return mux
}
