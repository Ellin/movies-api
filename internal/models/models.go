package models

type Genre struct {
	ID   int
	Name string
}

type Movie struct {
	ID          int64
	Title       string
	ReleaseYear int
	Duration    int
	GenreIDs    []int
	ActorIDs    []int
}

type Actor struct {
	ID        int
	Name      string
	BirthDate string
}
