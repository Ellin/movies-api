package models

type Genre struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	MovieIDs []int64 `json:"movie_ids"`
}

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
	Genres      []Genre        `json:"genres"`
	Actors      []ActorSummary `json:"actors"`
}

type Actor struct {
	ID        int64
	Name      string
	BirthDate string
}

type ActorSummary struct {
	ID   int64
	Name string
}
