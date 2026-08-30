package models

import "movies-api/internal/pagination"

// Movie is used to query the database
type Movie struct {
	ID          int64
	Title       string
	ReleaseYear int
	Duration    int
	GenreIDs    []int64
	ActorIDs    []int64
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

// MovieFilter stores data related to filtering and pagination from query parameters
type MovieFilter struct {
	ReleaseYear *int
	Genre       *int64
	Actor       *int64
	Pagination  pagination.Pagination
}
