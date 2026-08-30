package main

import (
	"movies-api/internal/handlers"
	"movies-api/internal/middleware"
	"net/http"
	"time"
)

func NewRouter(app *handlers.App) http.Handler {
	mux := http.NewServeMux()

	// movie routes
	mux.HandleFunc("POST /api/movie", app.PostMovie)
	mux.HandleFunc("GET /api/movie", app.GetAllMovies)
	mux.HandleFunc("GET /api/movie/{id}/actors", app.GetMovieActors)
	mux.HandleFunc("GET /api/movie/{id}", app.GetMovie)
	mux.HandleFunc("GET /api/movie/search", app.GetMovieSearch)
	mux.HandleFunc("PATCH /api/movie/{id}", app.PatchMovie)
	mux.HandleFunc("DELETE /api/movie/{id}", app.DeleteMovie)
	//genre routes
	mux.HandleFunc("POST /api/genre", app.PostGenre)
	mux.HandleFunc("GET /api/genre", app.GetAllGenres)
	mux.HandleFunc("GET /api/genre/{id}", app.GetGenre)
	mux.HandleFunc("PATCH /api/genre/{id}", app.PatchGenre)
	mux.HandleFunc("DELETE /api/genre/{id}", app.DeleteGenre)
	//actor routes
	mux.HandleFunc("POST /api/actor", app.PostActor)
	mux.HandleFunc("GET /api/actor", app.GetAllActors)
	mux.HandleFunc("GET /api/actor/{id}", app.GetActor)
	mux.HandleFunc("PATCH /api/actor/{id}", app.PatchActor)
	mux.HandleFunc("DELETE /api/actor/{id}", app.DeleteActor)

	// Add middleware
	var handler http.Handler = mux // This works because *http.ServeMux satisfies the http.Handler interface
	handler = middleware.Timeout(5 * time.Second)(handler)
	handler = middleware.RecoverPanic(handler)

	return handler
}
