package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// movie routes
	mux.HandleFunc("POST /api/movie", app.postMovie)
	mux.HandleFunc("GET /api/movie", app.getAllMovies)
	mux.HandleFunc("GET /api/movie/{id}", app.getMovie)

	return mux
}

// Implement HTTP handler functions for each entity (Genre, Movie, and Actor)
// Set up the following endpoints for each entity:

// POST /api/{entity}: Create a new entity
// GET /api/{entity}: Retrieve all entities
// GET /api/{entity}/{id}: Retrieve a specific entity by ID
// PATCH /api/{entity}/{id}: Partially update an existing entity
// DELETE /api/{entity}/{id}: Delete an entity
// Additionally, implement filtering endpoints for the following:

// GET /api/movies?genre={genreId}: Retrieve movies filtered by genre
// GET /api/movies?year={releaseYear}: Retrieve movies filtered by release year
// GET /api/movies?actor={actorId}: Retrieve movies that the actor with the given id has starred in
// GET /api/movies/{movieId}/actors: Retrieve all actors starring in a movie
// GET /api/actors?name={name}: Retrieve actors filtered by name
