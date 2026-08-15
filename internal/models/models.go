package models

type Genre struct {
	ID   int
	Name string
}

type Movie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	Duration    int    `json:"duration"`
	GenreIDs    []int  `json:"genre_ids"`
	ActorIDs    []int  `json:"actor_ids"`
}

type Actor struct {
	ID        int
	Name      string
	BirthDate string
}
