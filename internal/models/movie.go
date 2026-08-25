package models

import "movies-api/internal/pagination"

// movie structs
type Movie struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear int     `json:"release_year"`
	Duration    int     `json:"duration"`
	GenreIDs    []int64 `json:"genre_ids"`
	ActorIDs    []int64 `json:"actor_ids"`
}

// MovieDetail contains more detailed genre and actor info for JSON responses
type MovieDetail struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	ReleaseYear int            `json:"release_year"`
	Duration    int            `json:"duration"`
	Genres      []GenreSummary `json:"genres"`
	Actors      []ActorSummary `json:"actors"`
}

type MovieFilter struct {
	ReleaseYear *int
	Genre       *int64
	Actor       *int64
	Pagination  pagination.Pagination
}
