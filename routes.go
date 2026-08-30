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
	mux.HandleFunc("POST /api/movies", app.PostMovie)
	mux.HandleFunc("GET /api/movies", app.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}/actors", app.GetMovieActors)
	mux.HandleFunc("GET /api/movies/{id}", app.GetMovie)
	mux.HandleFunc("GET /api/movies/search", app.GetMovieSearch)
	mux.HandleFunc("PATCH /api/movies/{id}", app.PatchMovie)
	mux.HandleFunc("DELETE /api/movies/{id}", app.DeleteMovie)
	//genre routes
	mux.HandleFunc("POST /api/genres", app.PostGenre)
	mux.HandleFunc("GET /api/genres", app.GetAllGenres)
	mux.HandleFunc("GET /api/genres/{id}", app.GetGenre)
	mux.HandleFunc("PATCH /api/genres/{id}", app.PatchGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", app.DeleteGenre)
	//actor routes
	mux.HandleFunc("POST /api/actors", app.PostActor)
	mux.HandleFunc("GET /api/actors", app.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", app.GetActor)
	mux.HandleFunc("PATCH /api/actors/{id}", app.PatchActor)
	mux.HandleFunc("DELETE /api/actors/{id}", app.DeleteActor)

	// Add middleware
	var handler http.Handler = mux // This works because *http.ServeMux satisfies the http.Handler interface
	handler = middleware.Timeout(5 * time.Second)(handler)
	handler = middleware.RecoverPanic(handler)

	return handler
}
